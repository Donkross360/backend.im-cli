package commands

import (
	"fmt"

	"github.com/backend-im/cli/internal/api"
	"github.com/backend-im/cli/internal/auth"
	"github.com/backend-im/cli/internal/files"
	"github.com/spf13/cobra"
)

func NewCommitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit [project-id]",
		Short: "Commit local changes to Backend.im",
		Long:  "Commit local file changes to Backend.im. Changes are saved to Gitea and can be deployed.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, _ := cmd.Flags().GetString("dir")
			projectID, _ := cmd.Flags().GetString("project")
			message, _ := cmd.Flags().GetString("message")

			// Get project ID from positional argument or flag
			if projectID == "" && len(args) > 0 {
				projectID = args[0]
			}

			if projectDir == "" {
				projectDir = "."
			}

			if projectID == "" {
				return fmt.Errorf("project ID is required (use: commit <project-id> or --project flag)")
			}

			if message == "" {
				message = "Update code from CLI"
			}

			// Load auth token
			token, err := auth.LoadToken()
			if err != nil {
				return fmt.Errorf("authentication required: %w\nRun 'backend-im login' first", err)
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
			fmt.Printf("💬 Commit message: %s\n", message)

			// Create API client
			apiClient := api.NewClient()
			apiClient.SetAuthToken(token.AccessToken)

			// Commit changes
			fmt.Println("💾 Committing changes to Backend.im...")
			commitResp, err := apiClient.CommitChanges(fileMap, projectID, message)
			if err != nil {
				return fmt.Errorf("failed to commit changes: %w", err)
			}

			fmt.Printf("✅ Changes committed successfully!\n")
			fmt.Printf("🔑 Commit Hash: %s\n", commitResp.CommitHash)
			fmt.Printf("📊 Status: %s\n", commitResp.Status)
			fmt.Println("")
			fmt.Println("💡 Next step: Run 'backend-im deploy <project-id>' to deploy your changes")

			return nil
		},
	}

	cmd.Flags().StringP("dir", "d", "", "Project directory (default: current directory)")
	cmd.Flags().StringP("project", "p", "", "Project ID (can also be provided as positional argument)")
	cmd.Flags().StringP("message", "m", "", "Commit message (default: 'Update code from CLI')")

	return cmd
}

