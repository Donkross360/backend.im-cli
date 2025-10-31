package files

import (
	"fmt"
	"os"
	"path/filepath"
)

func DownloadFiles(files map[string]string, projectDir string) error {
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	for filename, content := range files {
		filePath := filepath.Join(projectDir, filename)

		// Create directory if needed
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", filename, err)
		}

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filename, err)
		}
	}

	return nil
}

func ReadProjectFiles(projectDir string) (map[string]string, error) {
	files := make(map[string]string)

	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Skip hidden files and common ignore patterns
		baseName := filepath.Base(path)
		if len(baseName) > 0 && baseName[0] == '.' {
			return nil
		}

		// Skip common ignore patterns
		relPath, _ := filepath.Rel(projectDir, path)
		if shouldIgnore(relPath) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		files[relPath] = string(content)
		return nil
	})

	return files, err
}

func shouldIgnore(path string) bool {
	ignorePatterns := []string{
		".git/",
		"__pycache__/",
		"*.pyc",
		".env",
		"node_modules/",
		".backend-im/",
	}

	for _, pattern := range ignorePatterns {
		matched, _ := filepath.Match(pattern, path)
		if matched {
			return true
		}
		// Check if path contains the pattern
		if len(pattern) > 0 && pattern[len(pattern)-1] == '/' {
			if len(path) >= len(pattern) && path[:len(pattern)] == pattern {
				return true
			}
		}
	}

	return false
}

