package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/backend-im/cli/internal/api"
	"github.com/backend-im/cli/internal/auth"
	"github.com/backend-im/cli/internal/files"
	"github.com/spf13/cobra"
)

func NewDeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy [project-id]",
		Short: "Deploy local code to Backend.im",
		Long:  "Deploy local code files to Backend.im. Backend.im will commit to Gitea automatically.",
		RunE: func(cmd *cobra.Command, args []string) error {
			watch, _ := cmd.Flags().GetBool("watch")
			projectDir, _ := cmd.Flags().GetString("dir")
			projectID, _ := cmd.Flags().GetString("project")

			// Get project ID from positional argument or flag
			if projectID == "" && len(args) > 0 {
				projectID = args[0]
			}

			if projectDir == "" {
				projectDir = "."
			}

			if projectID == "" {
				return fmt.Errorf("project ID is required (use: deploy <project-id> or --project flag)")
			}

			// Load auth token
			token, err := auth.LoadToken()
			if err != nil {
				return fmt.Errorf("authentication required: %w\nRun 'backend-im auth' first", err)
			}

			// Read local project files
			fmt.Printf("📂 Reading files from: %s\n", projectDir)
			fileMap, err := files.ReadProjectFiles(projectDir)
			if err != nil {
				return fmt.Errorf("failed to read project files: %w", err)
			}

			if len(fileMap) == 0 {
				return fmt.Errorf("no files found in %s", projectDir)
			}

			fmt.Printf("📦 Found %d files\n", len(fileMap))
			fmt.Printf("📁 Project ID: %s\n", projectID)

			// Create API client
			apiClient := api.NewClient()
			apiClient.SetAuthToken(token.AccessToken)

			// Deploy
			fmt.Println("🚀 Deploying to Backend.im...")
			deployResp, err := apiClient.Deploy(fileMap, projectID)
			if err != nil {
				return fmt.Errorf("deployment failed: %w", err)
			}

			fmt.Printf("✅ Deployment started!\n")
			fmt.Printf("📋 Deployment ID: %s\n", deployResp.DeploymentID)
			fmt.Printf("🔑 Commit Hash: %s\n", deployResp.CommitHash)
			fmt.Printf("📊 Status: %s\n", deployResp.Status)

			if watch {
				// Use WebSocket for real-time updates
				if err := streamDeploymentUpdates(apiClient, deployResp.DeploymentID, deployResp.WebSocketURL); err != nil {
					return fmt.Errorf("failed to stream updates: %w", err)
				}
			} else {
				// Poll for deployment status until we get the URL
				fmt.Println("⏳ Waiting for deployment to complete...")
				url, err := pollForDeploymentURL(apiClient, deployResp.DeploymentID)
				if err != nil {
					fmt.Fprintf(os.Stderr, "⚠️  Warning: Could not get deployment URL: %v\n", err)
					fmt.Println("💡 Use --watch flag to see real-time progress")
				} else if url != "" {
					fmt.Printf("🌐 Deployment URL: %s\n", url)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolP("watch", "w", true, "Watch deployment progress in real-time (default: true)")
	cmd.Flags().StringP("dir", "d", "", "Project directory (default: current directory)")
	cmd.Flags().StringP("project", "p", "", "Project ID (can also be provided as positional argument)")

	return cmd
}

// pollForDeploymentURL polls the status endpoint until deployment completes and returns URL
func pollForDeploymentURL(apiClient *api.Client, deploymentID string) (string, error) {
	maxAttempts := 30 // 30 attempts = ~30 seconds
	attempt := 0
	lastStatus := ""

	for attempt < maxAttempts {
		status, err := apiClient.GetStatus(deploymentID)
		if err != nil {
			return "", err
		}

		// Only show status when it changes
		if status.Status != lastStatus {
			fmt.Printf("📊 Status: %s", status.Status)
			if len(status.Logs) > 0 {
				fmt.Printf(" - %s", status.Logs[len(status.Logs)-1])
			}
			fmt.Println()
			lastStatus = status.Status
		}

		// If deployment is complete, return URL
		if status.Status == "complete" {
			if status.URL != "" {
				return status.URL, nil
			}
		}

		// If deployment failed, return error
		if status.Status == "failed" {
			return "", fmt.Errorf("deployment failed")
		}

		// Wait before next poll
		time.Sleep(1 * time.Second)
		attempt++
	}

	return "", fmt.Errorf("timeout waiting for deployment URL (after %d seconds)", maxAttempts)
}

// streamDeploymentUpdates connects to WebSocket and streams real-time deployment updates
func streamDeploymentUpdates(apiClient *api.Client, deploymentID, websocketURL string) error {
	// Always use base URL from API client for consistency (especially important in Docker)
	// The API client's base URL is correctly configured for the environment (e.g., http://mock-api:8080)
	// The websocketURL from the API response may contain localhost URLs that don't work in containerized environments
	baseURL := apiClient.BaseURL()

	wsClient := api.NewWebSocketClient(baseURL)
	defer wsClient.Close()

	fmt.Println("🔌 Connecting to WebSocket...")
	if err := wsClient.Connect(deploymentID); err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	fmt.Println("👀 Streaming deployment updates...")
	fmt.Println("")

	lastStatus := ""
	var finalURL string

	err := wsClient.StreamUpdates(func(update *api.DeploymentUpdate) error {
		// Show status when it changes
		if update.Status != lastStatus {
			fmt.Printf("📊 Status: %s", update.Status)
			if update.Namespace != "" {
				fmt.Printf(" (namespace: %s)", update.Namespace)
			}
			if update.PVC != "" {
				fmt.Printf(" (PVC: %s)", update.PVC)
			}
			fmt.Println()
			lastStatus = update.Status
		}

		// Show logs
		for _, log := range update.Logs {
			fmt.Printf("   %s\n", log)
		}

		// Store URL when available (may come before status="complete")
		if update.URL != "" {
			finalURL = update.URL
		}

		// Exit conditions
		if update.Status == "complete" {
			fmt.Println("")
			if finalURL != "" {
				fmt.Printf("🌐 Deployment URL: %s\n", finalURL)
			} else {
				fmt.Println("✅ Deployment completed successfully")
			}
			return nil
		}

		if update.Status == "failed" {
			return fmt.Errorf("deployment failed")
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("WebSocket streaming error: %w", err)
	}

	return nil
}

