package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"register-deployed-artifact-test/internal/deployments"

	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{
		Use:   "register-deployed-test-action",
		Short: "Publish the deployment metadata to CloudBees Platform",
		Long:  "Publish the deployment metadata to CloudBees Platform",
		RunE:  run,
	}
	cfg deployments.Config
)

func Execute() error {
	return cmd.Execute()
}

func init() {
	cmd.Flags().BoolVar(&cfg.DryRun, "dry-run", false,
		"Dry run mode - validate configuration without sending events to CloudBees Platform")
}

func run(_ *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("too many arguments: %v", args)
	}
	newContext, cancel := context.WithCancel(context.Background())
	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, os.Interrupt)
	go func() {
		<-osChannel
		cancel()
	}()

	return cfg.Run(newContext)
}
