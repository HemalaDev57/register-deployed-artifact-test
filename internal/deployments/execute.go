package deployments

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"

	github "github.com/calculi-corp/actions-common/github"
	commons "github.com/calculi-corp/actions-common/rest"
)

// To allow mocking in test, exported as it's mocked in root_test.go too
var SendCloudEvent = commons.SendCloudEvent

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
	}else {
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
	cfg.GhDetails = github.GetGithubEnvVars()

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
