package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const AppName = "web3safe"

func GetDefaultConfigDirectory() string {
	var configDir string

	switch osType := runtime.GOOS; osType {
	case "windows":
		appData := os.Getenv("APPDATA")
		configDir = filepath.Join(appData, AppName)
	case "darwin", "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic("Unable to get home directory")
		}

		configDir = filepath.Join(homeDir, ".config", AppName)
	default:
		panic("Unsupported operating system")
	}

	return configDir
}

func IsFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}

		panic(fmt.Sprintf("unable to check file existence: %v", err))
	}

	return true
}

func GetFileContent(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return strings.Split(string(content), "\n"), nil
}
