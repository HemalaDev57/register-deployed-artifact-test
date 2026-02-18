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

	// Token related env vars, will always be avilable if action has permission `id-token: write`
	ActionIDTokenRequestURL   = "ACTIONS_ID_TOKEN_REQUEST_URL"
	ActionIDTokenRequestToken = "ACTIONS_ID_TOKEN_REQUEST_TOKEN"
	// Github env vars
	GithubJob                  = "GITHUB_JOB"
	GithubOutput               = "GITHUB_OUTPUT"
	GithubRepository           = "GITHUB_REPOSITORY"
	GithubRunAttempt           = "GITHUB_RUN_ATTEMPT"
	GithubRunID                = "GITHUB_RUN_ID"
	GithubRunNumber            = "GITHUB_RUN_NUMBER"
	GithubURL                  = "GITHUB_SERVER_URL"
	GithubWorkflowRef          = "GITHUB_WORKFLOW_REF"
	AccessToken                = "accessToken"
	AuthorizationHeaderKey     = "Authorization"
	Bearer                     = "Bearer "
	ContentTypeHeaderKey       = "Content-Type"
	ContentTypeCloudEventsJSON = "application/cloudevents+json"
	GithubProvider             = "GITHUB"
)
