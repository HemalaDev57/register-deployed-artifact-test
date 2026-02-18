package deployments

import (
	"context"
	"fmt"
	"log"
	"os"
)

type Config struct {
	context.Context
	TargetEnvironment string
	DeploymentLabels  string
	ArtifactID        string
	ArtifactURL       string
	CloudBeesAPIURL   string
	DryRun            bool
	GhDetails         GithubDetails
}

type GithubDetails struct {
	GithubJob         string
	GithubRepository  string
	GithubRunAttempt  string
	GithubRunID       string
	GithubRunNumber   string
	GithubURL         string
	GithubWorkflowRef string
}

// Injects Github environment variables into a GithubDetails object.
// These env vars will always be available when running a GHA
func GetGithubEnvVars() GithubDetails {
	return GithubDetails{
		GithubJob:         os.Getenv(GithubJob),
		GithubRepository:  os.Getenv(GithubRepository),
		GithubRunAttempt:  os.Getenv(GithubRunAttempt),
		GithubRunID:       os.Getenv(GithubRunID),
		GithubRunNumber:   os.Getenv(GithubRunNumber),
		GithubURL:         os.Getenv(GithubURL),
		GithubWorkflowRef: os.Getenv(GithubWorkflowRef),
	}
}

// Writes the provided data into the file indicated by GITHUB_OUTPUT env var
// data will be written as `key=value\n` for each map entry
func WriteGitHubOutput(dataToWrite map[string]string) error {
	outputFilePath := os.Getenv(GithubOutput)
	// Open the GITHUB_OUTPUT file to append the output
	outputFile, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return fmt.Errorf("error opening GITHUB_OUTPUT (%s) file: %v", outputFilePath, err)
	}
	defer func(file os.File) {
		if err := file.Close(); err != nil {
			log.Printf("file %s was unexpctedly closed:%v", outputFilePath, err)
		}
	}(*outputFile)

	// Add all the provided keys and values to output file
	for k, v := range dataToWrite {
		_, err := fmt.Fprintf(outputFile, "%s=%s\n", k, v)
		if err != nil {
			return fmt.Errorf("error writing to GITHUB_OUTPUT (%s): %v", outputFilePath, err)
		}
	}
	return nil
}
