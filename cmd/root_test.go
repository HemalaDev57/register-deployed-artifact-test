// //nolint:errcheck
package cmd

// import (
// 	"gha-register-deployed-artifact/internal/deployments"
// 	"testing"

// 	v2 "github.com/cloudevents/sdk-go/v2"
// 	"github.com/stretchr/testify/assert"
// )

// func TestRunWithArtifactID(t *testing.T) {
// 	t.Setenv(deployments.TargetEnvironment, "production")
// 	t.Setenv(deployments.ArtifactID, "artifact-123")
// 	t.Setenv(deployments.ArtifactURL, "") // Clear artifact URL

// 	// Mocking response from actions-common
// 	SendCloudEvent = func(cloudEvent v2.Event, cloudbeesAPIURL string, endpointPath string) ([]byte, error) {
// 		return nil, nil
// 	}

// 	t.Setenv(deployments.CloudBeesAPIURL, "https://example.test")
// 	err := run(nil, nil)
// 	assert.Nil(t, err)
// }

// func TestRunWithArtifactURL(t *testing.T) {
// 	t.Setenv(deployments.TargetEnvironment, "staging")
// 	t.Setenv(deployments.ArtifactID, "") // Clear artifact ID
// 	t.Setenv(deployments.ArtifactURL, "https://docker.io/myapp:1.0.0")

// 	// Mocking response from actions-common
// 	deployments.SendCloudEvent = func(cloudEvent v2.Event, cloudbeesAPIURL string, endpointPath string) ([]byte, error) {
// 		return nil, nil
// 	}

// 	t.Setenv(deployments.CloudBeesAPIURL, "https://example.test")
// 	err := run(nil, nil)
// 	assert.Nil(t, err)
// }

// func TestFailureUnknownArguments(t *testing.T) {
// 	err := run(nil, []string{"test", "command"})
// 	assert.Contains(t, err.Error(), "too many arguments:")
// }

// func TestMissingArtifactIDAndUrl(t *testing.T) {
// 	t.Setenv(deployments.TargetEnvironment, "production")
// 	t.Setenv(deployments.ArtifactID, "")
// 	t.Setenv(deployments.ArtifactURL, "")

// 	t.Setenv(deployments.CloudBeesAPIURL, "https://example.test.com")
// 	err := run(nil, nil)
// 	assert.Contains(t, err.Error(), "either ARTIFACT_ID or ARTIFACT_URL must be provided")
// }
