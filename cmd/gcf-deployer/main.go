package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/api/cloudfunctions/v1"
)

var (
	version      string
	projectID    string
	locationID   string
	functionName string
	sourceURL    string
)

func deploy() error {
	ctx := context.Background()
	cloudfunctionsService, err := cloudfunctions.NewService(ctx)
	if err != nil {
		return err
	}

	projectsLocationsFunctionsService := cloudfunctions.NewProjectsLocationsFunctionsService(cloudfunctionsService)

	gcfName := fmt.Sprintf("projects/%s/locations/%s/functions/%s", projectID, locationID, functionName)
	cloudfunction, err := projectsLocationsFunctionsService.Get(gcfName).Context(ctx).Do()
	if err != nil {
		return err
	}

	// for debug
	fmt.Fprintf(os.Stderr, "%v\n", cloudfunction)
	fmt.Fprintf(os.Stderr, "%v\n", cloudfunction.SourceRepository.DeployedUrl)

	cloudfunction.SourceRepository.Url = sourceURL
	operation, err := projectsLocationsFunctionsService.Patch(gcfName, cloudfunction).Context(ctx).UpdateMask("sourceRepository.url").Do()
	if err != nil {
		return err
	}

	// wait for update function
	operationService := cloudfunctions.NewOperationsService(cloudfunctionsService)
	for {
		if err != nil {
			return err
		}
		if operation.Done {
			break
		}
		time.Sleep(5 * time.Second)
		fmt.Fprintf(os.Stderr, ".")
		operation, err = operationService.Get(operation.Name).Context(ctx).Do()
	}

	return nil
}

func main() {
	rootCmd := &cobra.Command{
		Use: "gcf-deployer",
		RunE: func(c *cobra.Command, args []string) error {
			return deploy()
		},
	}

	rootCmd.PersistentFlags().StringVar(&projectID, "project", "", "your project id")
	rootCmd.PersistentFlags().StringVar(&locationID, "region", "", "location of the function")
	rootCmd.PersistentFlags().StringVar(&functionName, "name", "", "your function name")
	rootCmd.PersistentFlags().StringVar(&sourceURL, "source", "", "location of source code to deploy")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
