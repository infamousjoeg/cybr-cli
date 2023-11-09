package util

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	terminal "golang.org/x/term"
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

// PromptForPassword prompts the user to enter a password, retrying until valid input is provided or max attempts are reached
func PromptForPassword(maxAttempts int) (string, error) {
	var password string
	var err error

	for attempts := 0; attempts < maxAttempts; attempts++ {
		password, err = ReadPassword()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if password != "" {
			return password, nil
		}
		fmt.Println("Password cannot be empty. Please try again.")
	}
	return "", fmt.Errorf("No valid password entered after %d attempts. Exiting", maxAttempts)
}

// ReadInput Read input from Stdin
func ReadInput(message string) (string, error) {
	fmt.Printf("%s: ", message)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", fmt.Errorf("An error occurred trying to read input from Stdin. Exiting")
	}
	return input, nil
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

// GetSubDomain Get the subdomain from the platform URL
func GetSubDomain(platformURL string) (string, error) {
	// Get the subdomain from the platform URL https://<subdomain>.privilegecloud.cyberark.cloud
	parsedURL, err := url.Parse(platformURL)
	if err != nil {
		return "", fmt.Errorf("Failed to parse URL. %s", err)
	}

	parts := strings.Split(parsedURL.Hostname(), ".")
	if len(parts) > 2 {
		return parts[0], nil
	}

	return "", fmt.Errorf("Failed to get subdomain from URL. %s", err)
}
