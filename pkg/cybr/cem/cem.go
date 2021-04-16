package cem

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/util"
)

// SaveToken saving token as file on the local filesystem
func SaveToken(token string, tokenPath string) error {
	// Get user home directory
	userHome, err := util.GetUserHomeDir()
	if err != nil {
		return fmt.Errorf("ACL error. %s", err)
	}

	// Check if .cybr directory already exists, create if not
	if _, err = os.Stat(userHome + "/.cybr"); os.IsNotExist(err) {
		// Create .cybr folder in user home directory
		err = os.Mkdir(userHome+"/.cybr", 0766)
		if err != nil {
			return fmt.Errorf("could not create folder %s/.cybr on local file system. %s", userHome, err)
		}
	}

	// Check for config file and remove if existing
	if _, err = os.Stat(userHome + tokenPath); !os.IsNotExist(err) {
		err = os.Remove(userHome + tokenPath)
		if err != nil {
			return fmt.Errorf("could not remove existing %s%s file. %s", userHome, tokenPath, err)
		}
	}
	// Create config file in user home directory
	dataFile, err := os.Create(userHome + tokenPath)
	if err != nil {
		return fmt.Errorf("could not create configuration file at %s%s. %s", userHome, tokenPath, err)
	}

	// serialize the data
	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(token)

	dataFile.Close()

	return nil
}

// GetToken file from local filesystem and read
func GetToken(tokenPath string) (string, error) {

	// Get user home directory
	userHome, err := util.GetUserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ACL error. %s", err)
	}

	// open data file
	dataFile, err := os.Open(userHome + tokenPath)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve token file at %s. %s", tokenPath, err)
	}

	dataDecoder := gob.NewDecoder(dataFile)
	result := ""
	err = dataDecoder.Decode(&result)
	if err != nil {
		return result, fmt.Errorf("failed to decode token file at .cybr/config. %s", err)
	}

	dataFile.Close()

	return result, nil
}
