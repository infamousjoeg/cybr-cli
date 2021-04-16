package util

import (
	"fmt"
	"os"
)

// GetUserHomeDir Get the Home directory of the current user
func GetUserHomeDir() (string, error) {
	// Get user home directory
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not read user home directory for OS. %s", err)
	}
	return userHome, nil
}
