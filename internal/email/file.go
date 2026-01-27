package email

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// openFile opens a file for reading
func openFile(filePath string) (*os.File, error) {
	return os.Open(filePath)
}

// CollectFiles returns all .eml files in the given directory
func CollectFiles(dirPath string) ([]string, error) {
	var emailFiles []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".eml") {
			emailFiles = append(emailFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return emailFiles, nil
}
