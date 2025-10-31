package editor

import (
	"fmt"
	"os/exec"
)

func OpenEditor(editorName, projectDir string) error {
	var cmd *exec.Cmd

	switch editorName {
	case "vscode", "code":
		cmd = exec.Command("code", projectDir)
	case "vim":
		cmd = exec.Command("vim", projectDir)
	case "nano":
		cmd = exec.Command("nano", projectDir)
	default:
		// Try to run as-is (user might specify full path)
		cmd = exec.Command(editorName, projectDir)
	}

	// Don't wait - let editor run independently
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open editor %s: %w", editorName, err)
	}

	return nil
}

