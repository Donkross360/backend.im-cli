package commands

import (
	"fmt"
	"os"

	"github.com/backend-im/cli/internal/auth"
	"github.com/backend-im/cli/internal/editor"
	"github.com/spf13/cobra"
)

func NewEditCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit [project-directory]",
		Short: "Open project directory in editor",
		Long:  "Open the project directory in your default editor to view and edit code.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, _ := cmd.Flags().GetString("dir")
			editorName, _ := cmd.Flags().GetString("editor")

			// Get directory from positional argument or flag
			if projectDir == "" && len(args) > 0 {
				projectDir = args[0]
			}

			if projectDir == "" {
				projectDir = "."
			}

			// Check if directory exists
			if _, err := os.Stat(projectDir); os.IsNotExist(err) {
				return fmt.Errorf("directory does not exist: %s", projectDir)
			}

			// Load auth token (optional - just for verification)
			_, err := auth.LoadToken()
			if err != nil {
				return fmt.Errorf("authentication required: %w\nRun 'backend-im login' first", err)
			}

			// Determine editor to use
			if editorName == "" {
				editorName = os.Getenv("EDITOR")
				if editorName == "" {
					// Try common editors
					for _, e := range []string{"code", "vscode", "vim", "nano"} {
						if _, err := os.Stat("/usr/bin/" + e); err == nil {
							editorName = e
							break
						}
					}
					if editorName == "" {
						return fmt.Errorf("no editor specified and $EDITOR not set. Use --editor flag or set $EDITOR environment variable")
					}
				}
			}

			fmt.Printf("🔧 Opening %s in %s...\n", projectDir, editorName)
			if err := editor.OpenEditor(editorName, projectDir); err != nil {
				return fmt.Errorf("failed to open editor: %w", err)
			}

			fmt.Println("")
			fmt.Println("💡 After editing, use 'backend-im commit <project-id>' to save your changes")

			return nil
		},
	}

	cmd.Flags().StringP("dir", "d", "", "Project directory (default: current directory)")
	cmd.Flags().StringP("editor", "e", "", "Editor to use (default: $EDITOR or auto-detect)")

	return cmd
}

