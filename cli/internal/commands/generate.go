package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/backend-im/cli/internal/api"
	"github.com/backend-im/cli/internal/auth"
	"github.com/backend-im/cli/internal/editor"
	"github.com/backend-im/cli/internal/files"
	"github.com/spf13/cobra"
)

func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [prompt]",
		Short: "Generate backend code from prompt via Backend.im API",
		Long:  "Generate FastAPI code from a prompt. Backend.im commits to Gitea automatically.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prompt := args[0]
			projectID, _ := cmd.Flags().GetString("project")
			outputDir, _ := cmd.Flags().GetString("output")
			editorName, _ := cmd.Flags().GetString("editor")

			// Load auth token
			token, err := auth.LoadToken()
			if err != nil {
				return fmt.Errorf("authentication required: %w\nRun 'backend-im auth' first", err)
			}

			// Create API client
			apiClient := api.NewClient()
			apiClient.SetAuthToken(token.AccessToken)

			fmt.Printf("🚀 Generating code from prompt: %s\n", prompt)
			if projectID != "" {
				fmt.Printf("📁 Project ID: %s\n", projectID)
			}

			// Call Backend.im API
			generatedFiles, err := apiClient.GenerateCode(prompt)
			if err != nil {
				return fmt.Errorf("failed to generate code: %w", err)
			}

			// Determine output directory
			if outputDir == "" {
				outputDir = fmt.Sprintf("backend-%d", time.Now().Unix())
			}

			// Download files locally
			fmt.Printf("📥 Downloading files to: %s\n", outputDir)
			if err := files.DownloadFiles(generatedFiles, outputDir); err != nil {
				return fmt.Errorf("failed to download files: %w", err)
			}

			fmt.Printf("✅ Code generated successfully in ./%s\n", outputDir)
			fmt.Printf("📝 Files: %d files downloaded\n", len(generatedFiles))
			fmt.Println("")
			fmt.Println("💡 The code has been committed to Backend.im/Gitea automatically.")
			fmt.Println("   You can now edit the files locally and use 'backend-im commit <project-id>' to save changes.")

			// Auto-open editor if specified, or try to detect default editor
			if editorName != "" {
				fmt.Printf("🔧 Opening %s...\n", editorName)
				if err := editor.OpenEditor(editorName, outputDir); err != nil {
					fmt.Fprintf(os.Stderr, "⚠️  Warning: Failed to open editor: %v\n", err)
				}
			} else {
				// Try to detect and open default editor from environment
				if defaultEditor := os.Getenv("EDITOR"); defaultEditor != "" {
					fmt.Printf("🔧 Opening %s...\n", defaultEditor)
					if err := editor.OpenEditor(defaultEditor, outputDir); err != nil {
						// Silent fail - not critical
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringP("project", "p", "", "Project ID (created when user creates project)")
	cmd.Flags().StringP("output", "o", "", "Output directory (default: auto-generated name)")
	cmd.Flags().StringP("editor", "e", "", "Open generated code in editor (vscode, code, vim, etc.). Auto-detects $EDITOR if not specified.")

	return cmd
}

