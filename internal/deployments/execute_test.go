package deployments

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"strings"
// 	"testing"
// 	"time"

// 	github "github.com/calculi-corp/actions-common/github"
// 	common "github.com/calculi-corp/actions-common/rest"
// 	cloudevents "github.com/cloudevents/sdk-go/v2"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// )

// const artifactID = "artifact-id"

// func TestCloudEventCreation(t *testing.T) {
// 	t.Run("Base event fields creation test", func(t *testing.T) {
// 		var ghDetails = github.GithubDetails{
// 			GithubRepository:  "test-org/test-repo",
// 			GithubRunAttempt:  "1",
// 			GithubRunID:       "5",
// 			GithubRunNumber:   "2",
// 			GithubURL:         "https://github.com",
// 			GithubWorkflowRef: "calculi-corp/workflow-template-gha-actions/.github/workflows/gha-action-testing.yml@refs/heads/main",
// 		}
// 		var config = Config{
// 			Context:   context.Background(),
// 			GhDetails: ghDetails,
// 		}
// 		cloudEvent := createCloudEvent(&config)
// 		assert.Equal(t, "1.0", cloudEvent.SpecVersion(), "specVersion is 1.0")
// 		assert.Equal(t, "application/json", cloudEvent.DataContentType(), "expecting right content type")

// 		err := uuid.Validate(cloudEvent.ID())
// 		assert.Nil(t, err, "UUID should be valid")

// 		assert.Equal(t, "https://github.com/test-org/test-repo", cloudEvent.Source(), "source is url + / + repo")
// 		assert.Equal(t, "cloudbees.platform.register.deployed.artifact", cloudEvent.Type(), "type is expected")

// 		subjectParts := strings.Split(cloudEvent.Subject(), "|")
// 		assert.Equal(t, ghDetails.GithubWorkflowRef, subjectParts[0], "Subject[0] = workflow ref")
// 		assert.Equal(t, ghDetails.GithubRunID, subjectParts[1], "Subject[1] = run id")
// 		assert.Equal(t, ghDetails.GithubRunAttempt, subjectParts[2], "Subject[2] = run attempt")
// 		assert.Equal(t, ghDetails.GithubRunNumber, subjectParts[3], "Subject[3] = run number")

// 		// Dev memo: This should always be strictly greater, adding orEqual just to avoid potential flakiness
// 		assert.GreaterOrEqual(t, time.Now(), cloudEvent.Time(), "We should have a past creation time")
// 	})
// 	t.Run("Deployed artifact fields test", func(t *testing.T) {
// 		var config = Config{
// 			Context:           context.Background(),
// 			ArtifactID:        artifactID,
// 			ArtifactURL:       "https://some-url",
// 			TargetEnvironment: "test-env",
// 			DeploymentLabels:  "label-1,label-2",
// 			GhDetails: github.GithubDetails{
// 				GithubRunID:      "100",
// 				GithubRunAttempt: "2",
// 				GithubRunNumber:  "1",
// 				GithubJob:        "test-job",
// 			},
// 		}
// 		content := populateEventData(&config)
// 		providerInfo := content.ProviderInfo
// 		assert.Equal(t, "100", providerInfo.RunID, "run_id should be 100")
// 		assert.Equal(t, "2", providerInfo.RunAttempt, "run_attempt should be 2")
// 		assert.Equal(t, "1", providerInfo.RunNumber, "run_number should be 1")
// 		assert.Equal(t, "test-job", providerInfo.JobName, "job_name should be test-job")
// 		assert.Equal(t, "GITHUB", providerInfo.Provider, "provider should be GITHUB")

// 		artifactInfo := content.ArtifactInfo

// 		assert.Equal(t, artifactID, artifactInfo.ArtifactID, "artifact_id should be artifact-id")
// 		assert.Equal(t, "https://some-url", artifactInfo.ArtifactURL, "artifact_id should be the indicated")
// 		assert.Equal(t, "test-env", artifactInfo.TargetEnvironment, "target_environment should be test-env")
// 		assert.Equal(t, "label-1,label-2", artifactInfo.ArtifactLabel, "artifact_label should be equal to deployment labels")

// 		// Payload plain text testing, will only be checking field names, values already checked
// 		cloudEvent := createCloudEvent(&config)
// 		_ = cloudEvent.SetData(*cloudevents.StringOfApplicationJSON(), content)
// 		jsonText := common.PrettyPrint(cloudEvent)
// 		plainJSON := make(map[string]interface{})
// 		_ = json.Unmarshal([]byte(jsonText), &plainJSON)

// 		assert.Contains(t, plainJSON, "specversion")
// 		assert.Contains(t, plainJSON, "id")
// 		assert.Contains(t, plainJSON, "source")
// 		assert.Contains(t, plainJSON, "type")
// 		assert.Contains(t, plainJSON, "subject")
// 		assert.Contains(t, plainJSON, "datacontenttype")
// 		assert.Contains(t, plainJSON, "time")
// 		assert.Contains(t, plainJSON, "data")

// 		data := plainJSON["data"].(map[string]interface{})
// 		provInfo := data["provider_info"].(map[string]interface{})
// 		assert.Contains(t, provInfo, "run_id")
// 		assert.Contains(t, provInfo, "run_attempt")
// 		assert.Contains(t, provInfo, "run_number")
// 		assert.Contains(t, provInfo, "job_name")
// 		assert.Contains(t, provInfo, "provider")

// 		artInfo := data["artifact_info"].(map[string]interface{})
// 		assert.Contains(t, artInfo, "artifact_id")
// 		assert.Contains(t, artInfo, "artifact_url")
// 		assert.Contains(t, artInfo, "target_environment")
// 		assert.Contains(t, artInfo, "artifact_label")
// 	})
// }

// func TestEnvVars(t *testing.T) {
// 	t.Run("No env vars set", func(t *testing.T) {
// 		config := Config{}
// 		err := setEnvVars(&config)
// 		assert.NotNil(t, err, "no env vars set")
// 	})
// 	t.Run("Missing mandatory env var", func(t *testing.T) {
// 		t.Setenv(ArtifactID, artifactID)
// 		config := Config{}
// 		err := setEnvVars(&config)
// 		assert.NotNil(t, err, "mandatory vars not present")
// 		assert.Contains(t, err.Error(), TargetEnvironment, "missing mandatory "+TargetEnvironment+" var")
// 	})
// 	t.Run("Missing both of the dependant vars", func(t *testing.T) {
// 		t.Setenv(TargetEnvironment, "env")
// 		config := Config{}
// 		err := setEnvVars(&config)
// 		assert.NotNil(t, err, "we must have either ID or URL defined")
// 		assert.Contains(t, err.Error(), ArtifactID, "missing either "+ArtifactID+" or "+ArtifactURL)
// 		assert.Contains(t, err.Error(), ArtifactURL, "missing either "+ArtifactID+" or "+ArtifactURL)
// 	})
// 	t.Run("All mandatory vars set", func(t *testing.T) {
// 		t.Setenv(TargetEnvironment, "env")
// 		t.Setenv(ArtifactID, artifactID)
// 		config := Config{}
// 		err := setEnvVars(&config)
// 		assert.Nil(t, err, "we must not have an error")
// 		assert.Equal(t, artifactID, config.ArtifactID, "Artifact Id is set")
// 		assert.Equal(t, "env", config.TargetEnvironment, "Target environment is set")
// 		assert.Equal(t, DefaultAPIURL, config.CloudBeesAPIURL, "API URL has default value")
// 		assert.Empty(t, config.DeploymentLabels, "not defined vars")
// 		assert.Empty(t, config.ArtifactURL, "not defined vars")
// 	})
// 	t.Run("All vars set", func(t *testing.T) {
// 		t.Setenv(TargetEnvironment, "target-env")
// 		t.Setenv(ArtifactURL, "artifact-url")
// 		t.Setenv(DeploymentLabels, "label-1, label-2")
// 		t.Setenv(CloudBeesAPIURL, "https://example.com")
// 		config := Config{}
// 		err := setEnvVars(&config)
// 		assert.Nil(t, err, "we must not have an error")
// 		assert.Equal(t, "artifact-url", config.ArtifactURL, "Artifact URL is set")
// 		assert.Equal(t, "target-env", config.TargetEnvironment, "Target environment is set")
// 		assert.Equal(t, "https://example.com", config.CloudBeesAPIURL, "API URL is set")
// 		assert.Equal(t, "label-1, label-2", config.DeploymentLabels, "labels are set")
// 		assert.Empty(t, config.ArtifactID, "not defined var") // Testing that any of bot ID and URL is needed
// 	})
// }

// func TestActionRun(t *testing.T) {
// 	setTestEnvVars := func(t *testing.T) {
// 		t.Setenv(TargetEnvironment, "env")
// 		t.Setenv(ArtifactID, artifactID)
// 	}
// 	t.Run("Execution fails if no env vars", func(t *testing.T) {
// 		config := Config{}
// 		err := config.Run(context.Background())
// 		assert.NotNil(t, err, "Execution fails because of env vars")
// 		assert.Contains(t, err.Error(), "TARGET_ENVIRONMENT")
// 	})
// 	t.Run("Execution fails because of platform error response", func(t *testing.T) {
// 		setTestEnvVars(t)
// 		config := Config{}
// 		// mocking an error response from platform
// 		SendCloudEvent = func(cloudEvent cloudevents.Event, cloudbeesAPIURL string, endpointPath string) ([]byte, error) {
// 			return nil, fmt.Errorf("error sending CloudEvent to platform")
// 		}
// 		err := config.Run(context.Background())
// 		assert.NotNil(t, err, "Execution fails because error response")
// 		assert.Contains(t, err.Error(), "platform")
// 	})
// 	t.Run("Execution success", func(t *testing.T) {
// 		setTestEnvVars(t)
// 		config := Config{}
// 		// mocking an error response from platform
// 		SendCloudEvent = func(cloudEvent cloudevents.Event, cloudbeesAPIURL string, endpointPath string) ([]byte, error) {
// 			return []byte(`Success`), nil
// 		}
// 		err := config.Run(context.Background())
// 		assert.Nil(t, err, "Execution doesn't fail")
// 	})
// }
