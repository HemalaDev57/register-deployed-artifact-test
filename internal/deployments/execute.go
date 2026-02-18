package deployments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

func (config *Config) Run(_ context.Context) (err error) {
	err = setEnvVars(config)
	if err != nil {
		return err
	}

	eventToSend := createCloudEvent(config)
	err = eventToSend.SetData(*cloudevents.StringOfApplicationJSON(), populateEventData(config))
	if err != nil {
		return err
	}

	if config.DryRun {
		fmt.Printf("Running in dry run mode, skipping sending event\n%s", eventToSend)
	} else {
		_, err = SendCloudEvent(eventToSend, config.CloudBeesAPIURL, EndpointPath)
		if err != nil {
			return fmt.Errorf("error in the platform request: %v", err)
		}
	}

	// Output doesn't need to be parsed, just print success message
	fmt.Printf("Deployed artifact event (%s) succesfully sent to Platform", eventToSend.ID())
	return nil
}

func setEnvVars(cfg *Config) error {
	// mandatory field
	targetEnvironment := os.Getenv(TargetEnvironment)
	if targetEnvironment == "" {
		return fmt.Errorf(TargetEnvironment + " needs to be provided")
	}
	cfg.TargetEnvironment = targetEnvironment

	// At least one of both fields must be present, backend will choose priority if both are
	artifactID := os.Getenv(ArtifactID)
	artifactURL := os.Getenv(ArtifactURL)
	if artifactID == "" && artifactURL == "" {
		return fmt.Errorf("either " + ArtifactID + " or " + ArtifactURL + " must be provided")
	}
	cfg.ArtifactID = artifactID
	cfg.ArtifactURL = artifactURL

	// Optional fields
	cloudBeesAPIURL := os.Getenv(CloudBeesAPIURL)
	if cloudBeesAPIURL == "" {
		cloudBeesAPIURL = DefaultAPIURL
	}
	cfg.CloudBeesAPIURL = cloudBeesAPIURL
	cfg.DeploymentLabels = os.Getenv(DeploymentLabels)

	// Github environment variables, always available
	cfg.GhDetails = GetGithubEnvVars()

	return nil
}

func getSubject(config *Config) string {
	return strings.Join([]string{
		config.GhDetails.GithubWorkflowRef,
		config.GhDetails.GithubRunID,
		config.GhDetails.GithubRunAttempt,
		config.GhDetails.GithubRunNumber,
	}, "|")
}

func createCloudEvent(config *Config) cloudevents.Event {
	cloudEvent := cloudevents.NewEvent()
	cloudEvent.SetSpecVersion(SpecVersion)
	cloudEvent.SetID(uuid.NewString())
	cloudEvent.SetSource(config.GhDetails.GithubURL + "/" + config.GhDetails.GithubRepository)
	cloudEvent.SetType(Type)
	cloudEvent.SetSubject(getSubject(config))
	cloudEvent.SetDataContentType(*cloudevents.StringOfApplicationJSON())
	cloudEvent.SetTime(time.Now())
	return cloudEvent
}

func populateEventData(config *Config) Content {
	artifactInfo := &ArtifactInfo{
		ArtifactID:        config.ArtifactID,
		ArtifactURL:       config.ArtifactURL,
		TargetEnvironment: config.TargetEnvironment,
		ArtifactLabel:     config.DeploymentLabels, // note that in the payload deployment labels are artifact label
	}
	providerInfo := &ProviderInfo{
		RunID:      config.GhDetails.GithubRunID,
		RunAttempt: config.GhDetails.GithubRunAttempt,
		RunNumber:  config.GhDetails.GithubRunNumber,
		JobName:    config.GhDetails.GithubJob,
		Provider:   Provider,
	}
	return Content{
		ProviderInfo: *providerInfo,
		ArtifactInfo: *artifactInfo,
	}
}

// Sends a cloudevent to Platform via HTTP call. Needs platform base URL and path.
// This method will get the OIDC token, use it to gather Patform JWT token and then send the event using it.
func SendCloudEvent(cloudEvent cloudevents.Event, cloudbeesAPIURL string, endpointPath string) ([]byte, error) {
	accessToken, err := getAccessToken(cloudbeesAPIURL)
	if err != nil {
		return nil, err
	}
	log.Println("Initiated sending the cloudEvent to platform")
	eventJSON, err := json.Marshal(cloudEvent)
	if err != nil {
		return nil, fmt.Errorf("error encoding CloudEvent JSON %s", err)
	}
	log.Println(PrettyPrint(cloudEvent))
	eventReq, err := http.NewRequest(http.MethodPost, getExternalEventlURL(cloudbeesAPIURL, endpointPath), bytes.NewBuffer(eventJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create event request: %w", err)
	}
	eventReq.Header.Set(ContentTypeHeaderKey, ContentTypeCloudEventsJSON)
	eventReq.Header.Set(AuthorizationHeaderKey, Bearer+accessToken)
	client := &http.Client{}
	eventResp, err := client.Do(eventReq)
	if err != nil {
		return nil, fmt.Errorf("error sending external event: %w", err)
	}

	eventBodyBytes, err := io.ReadAll(eventResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	defer func(body io.ReadCloser) {
		err := eventResp.Body.Close()
		if err != nil {
			log.Fatalf("Error when closing response body:%v", err)
		}
	}(eventResp.Body)

	if eventResp.StatusCode != http.StatusOK {
		bodyObj := ErrorResponse{}
		msg := string(eventBodyBytes)
		if err := json.Unmarshal(eventBodyBytes, &bodyObj); err == nil && bodyObj.Message != "" {
			msg = bodyObj.Message
		}
		return nil, fmt.Errorf("error sending CloudEvent to platform - %s : %s", eventResp.Status, msg)
	}

	return eventBodyBytes, nil
}

func getAccessToken(cloudbeesAPIURL string) (string, error) {
	log.Println("Fetching OIDC token")
	oidcToken, err := getOIDCToken(cloudbeesAPIURL)
	if err != nil {
		return "", fmt.Errorf("failed to get OIDC token - %s", err.Error())
	}

	log.Println("Initiated exchanging the OIDC Token with CBP token...")
	tokenRequestObj := TokenRequest{
		Provider: GithubProvider,
		Audience: strings.TrimSuffix(cloudbeesAPIURL, "/"), // Optional: omit or override
	}
	tokenReqJSON, err := json.Marshal(tokenRequestObj)
	if err != nil {
		return "", fmt.Errorf("error encoding CloudEvent JSON %s", err)
	}
	tokenReq, _ := http.NewRequest(http.MethodPost, getExternalTokenExchangeURL(cloudbeesAPIURL), bytes.NewBuffer(tokenReqJSON))
	tokenReq.Header.Set(ContentTypeHeaderKey, ContentTypeCloudEventsJSON)
	tokenReq.Header.Set(AuthorizationHeaderKey, Bearer+oidcToken)

	client := &http.Client{}
	tokenResp, err := client.Do(tokenReq)
	if err != nil {
		return "", fmt.Errorf("error exchanging token with platform - %s", err.Error())
	}

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			log.Fatalf("Error closing response body:%v", err)
		}
	}(tokenResp.Body)

	bodyBytes, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if tokenResp.StatusCode != http.StatusOK {
		bodyObj := ErrorResponse{}
		msg := string(bodyBytes)
		if err := json.Unmarshal(bodyBytes, &bodyObj); err == nil && bodyObj.Message != "" {
			msg = bodyObj.Message
		}
		return "", fmt.Errorf("error during token exchange - %s : %s", tokenResp.Status, msg)
	}

	var respMap map[string]any
	if err := json.Unmarshal(bodyBytes, &respMap); err != nil {
		return "", fmt.Errorf("failed to parse token exchange response: %w", err)
	}

	accessToken, ok := respMap[AccessToken].(string)
	if !ok || accessToken == "" {
		return "", fmt.Errorf("accessToken missing or invalid in response")
	}
	log.Println("Token exchange successful!")
	return accessToken, nil
}

func getOIDCToken(cloudbeesAPIURL string) (string, error) {
	oidcToken := os.Getenv(ActionIDTokenRequestToken)
	oidcBaseURL := os.Getenv(ActionIDTokenRequestURL)
	if oidcToken == "" || oidcBaseURL == "" {
		return "", fmt.Errorf("needed environment variables %s and %s not set", ActionIDTokenRequestToken, ActionIDTokenRequestURL)

	}
	oidcAudience := url.QueryEscape(strings.TrimSuffix(cloudbeesAPIURL, "/"))
	oidcURL := fmt.Sprintf("%s?audience=%s", oidcBaseURL, oidcAudience)

	oidcTokenReq, err := http.NewRequest("GET", oidcURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create OIDC request: %v", err)
	}
	oidcTokenReq.Header.Add(AuthorizationHeaderKey, Bearer+oidcToken)
	client := &http.Client{}
	oidcTokenResp, err := client.Do(oidcTokenReq)
	if err != nil {
		return "", fmt.Errorf("failed to execute OIDC request: %v", err)
	}
	defer func(body io.ReadCloser) {
		if err := oidcTokenResp.Body.Close(); err != nil {
			log.Fatalf("Error closing response body:%v", err)
		}
	}(oidcTokenResp.Body)

	if oidcTokenResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(oidcTokenResp.Body)
		return "", fmt.Errorf("OIDC token request failed. Status: %d, Body: %s", oidcTokenResp.StatusCode, string(body))
	}

	var oidcResp struct{ Value string }
	if err := json.NewDecoder(oidcTokenResp.Body).Decode(&oidcResp); err != nil {
		return "", fmt.Errorf("failed to decode OIDC response: %v", err)
	}
	if oidcResp.Value == "" {
		return "", fmt.Errorf("OIDC token value is empty")
	}
	return oidcResp.Value, nil
}

func getExternalTokenExchangeURL(url string) string {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url + "token-exchange/external-oidc-id-token"
}

func getExternalEventlURL(base string, path string) string {
	if !strings.HasSuffix(path, "/") {
		base += "/"
	}
	return base + strings.TrimPrefix(path, "/")
}

// PrettyPrint converts the input to JSON string
func PrettyPrint(in any) string {
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		log.Fatalf("error marshalling response: %v\n", err)
	}
	return string(data)
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details []any  `json:"details"`
}

type TokenRequest struct {
	Provider string `json:"provider"`
	Audience string `json:"audience"`
}
