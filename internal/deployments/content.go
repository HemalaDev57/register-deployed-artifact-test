package deployments

// Content of the payload to send inside the "data" field
type Content struct {
	ProviderInfo ProviderInfo `json:"provider_info"`
	ArtifactInfo ArtifactInfo `json:"artifact_info"`
}

type ProviderInfo struct {
	RunID      string `json:"run_id"`
	RunAttempt string `json:"run_attempt"`
	RunNumber  string `json:"run_number"`
	JobName    string `json:"job_name"`
	Provider   string `json:"provider"`
}

type ArtifactInfo struct {
	ArtifactID        string `json:"artifact_id,omitempty"`
	ArtifactURL       string `json:"artifact_url,omitempty"`
	TargetEnvironment string `json:"target_environment"`
	ArtifactLabel     string `json:"artifact_label"`
}
