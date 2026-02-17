package deployments

import (
	"context"

	github "github.com/calculi-corp/actions-common/github"
)

type Config struct {
	context.Context
	TargetEnvironment string
	DeploymentLabels  string
	ArtifactID        string
	ArtifactURL       string
	CloudBeesAPIURL   string
	DryRun            bool
	GhDetails         github.GithubDetails
}
