//go:build darwin || linux

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func ConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}
	return filepath.Join(homeDir, ".config", "gowatch"), nil
}
