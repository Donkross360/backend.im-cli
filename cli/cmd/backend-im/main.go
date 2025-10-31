package main

import (
	"fmt"
	"os"

	"github.com/backend-im/cli/internal/commands"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "backend-im",
		Short: "Backend.im CLI for seamless deployment",
		Long:  "A CLI tool for deploying backend code to Backend.im platform",
	}

	// Authentication
	rootCmd.AddCommand(commands.NewAuthCommand())  // Keep for backward compatibility
	rootCmd.AddCommand(commands.NewLoginCommand()) // Preferred command name
	
	// Code generation and editing
	rootCmd.AddCommand(commands.NewGenerateCommand())
	rootCmd.AddCommand(commands.NewEditCommand())
	
	// Git operations
	rootCmd.AddCommand(commands.NewCommitCommand())
	
	// Deployment
	rootCmd.AddCommand(commands.NewDeployCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

