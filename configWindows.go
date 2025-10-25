//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func ConfigFilePath() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return "", fmt.Errorf("APPDATA environment variable is not set")
	}
	return filepath.Join(appData, "gowatch"), nil
}
