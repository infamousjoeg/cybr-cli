package util

import (
	"fmt"
	"os"
	"syscall"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"golang.org/x/crypto/ssh/terminal"
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

// ReadPassword Read password from Stdin
func ReadPassword() (string, error) {
	fmt.Print("Enter password: ")
	byteSecretVal, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("An error occurred trying to read password from Stdin. Exiting")
	}
	return string(byteSecretVal), nil
}

// ReadOTPcode Read one-time passcode from Stdin
func ReadOTPcode(credentials requests.Logon) (requests.Logon, error) {
	fmt.Print("Enter one-time passcode: ")
	byteOTPCode, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return credentials, fmt.Errorf("An error occurred trying to read one-time passcode from Stdin. Exiting")
	}
	credentials.Password = string(byteOTPCode)
	fmt.Println()
	return credentials, nil
}
