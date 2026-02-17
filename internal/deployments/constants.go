package deployments

const (
	// Args names
	TargetEnvironment = "TARGET_ENVIRONMENT"
	DeploymentLabels  = "DEPLOYMENT_LABELS"
	ArtifactID        = "ARTIFACT_ID"
	ArtifactURL       = "ARTIFACT_URL"
	CloudBeesAPIURL   = "CLOUDBEES_API_URL"

	// Default values
	DefaultAPIURL = "https://api.cloudbees.io/"
	EndpointPath  = "v3/external-events"

	// Payload constants
	Provider    = "GITHUB"
	SpecVersion = "1.0"
	Type        = "cloudbees.platform.register.deployed.artifact"
)
