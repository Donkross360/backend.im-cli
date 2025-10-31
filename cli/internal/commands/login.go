package commands

import (
	"fmt"
	"os"

	"github.com/backend-im/cli/internal/api"
	"github.com/backend-im/cli/internal/auth"
	"github.com/spf13/cobra"
)

// NewLoginCommand creates a login command (alias for auth)
func NewLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Backend.im",
		Long:  "Login to Backend.im using Google OAuth. Required for first-time setup.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if already authenticated
			existingToken, err := auth.LoadToken()
			if err == nil && existingToken != nil {
				fmt.Println("✅ Already logged in!")
				
				// Verify token is still valid
				apiClient := api.NewClient()
				apiClient.SetAuthToken(existingToken.AccessToken)
				verifyResp, err := apiClient.VerifyAuth()
				if err == nil && verifyResp.Valid {
					fmt.Printf("👤 User: %s (%s)\n", verifyResp.UserID, verifyResp.Email)
					return nil
				}
				fmt.Println("⚠️  Token expired, please login again...")
			}

			fmt.Println("🔐 Logging in to Backend.im...")
			fmt.Println("")
			fmt.Println("For now, using mock authentication.")
			fmt.Println("In production, this will open a browser for Google OAuth.")
			fmt.Println("")

			// Mock authentication flow (for testing with mock API)
			// In production, this would:
			// 1. Open browser with OAuth URL
			// 2. Handle callback
			// 3. Exchange code for token
			
			// For mock API, we'll use a simple mock token
			mockToken := &auth.Token{
				AccessToken: "mock_token_" + fmt.Sprintf("%d", os.Getpid()),
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			}

			if err := auth.SaveToken(mockToken); err != nil {
				return fmt.Errorf("failed to save token: %w", err)
			}

			fmt.Println("✅ Login successful!")
			fmt.Printf("💾 Token saved to: ~/.backend-im/token.json\n")

			// Verify the token works
			apiClient := api.NewClient()
			apiClient.SetAuthToken(mockToken.AccessToken)
			verifyResp, err := apiClient.VerifyAuth()
			if err == nil && verifyResp.Valid {
				fmt.Printf("👤 User: %s (%s)\n", verifyResp.UserID, verifyResp.Email)
			}

			return nil
		},
	}

	return cmd
}

