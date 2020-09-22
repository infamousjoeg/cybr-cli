package api

import (
	"fmt"
	"log"
	"strings"

	httpJson "github.com/infamousjoeg/pas-api-go/pkg/helpers"
)

// Logon to PAS REST API Web Service
// Because we're using concurrentSession capability, this is only supported
// on PAS REST API v11.3 and above
func Logon(hostname string, username string, password string, authType string, concurrent bool) (string, error) {
	// Return error if unsupported authentication type chosen
	if authType != "cyberark" && authType != "ldap" {
		log.Fatal("Unsupported auth type used. Only 'cyberark' or 'ldap' supported")
	}

	// Set URL for request
	url := fmt.Sprintf("%s/PasswordVault/api/auth/%s/logon", hostname, authType)
	body := fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\", \"concurrentSession\": \"%t\"}", username, password, concurrent)

	// Send request and received response
	response, err := httpJson.SendRequestRaw(url, "POST", "", body)
	if err != nil {
		return "", fmt.Errorf("Failed to authenticate to the PAS REST API. %s", err)
	}
	// Trim "..." from response and return session token
	return strings.Trim(string(response), "\""), err
}

// Logoff the PAS REST API Web Service
func Logoff(hostname string, token string) (bool, error) {
	// Set URL for request
	url := fmt.Sprintf("%s/PasswordVault/api/auth/logoff", hostname)
	_, err := httpJson.Post(url, token, "")
	if err != nil {
		return false, fmt.Errorf("Unable to logoff PAS REST API Web Service. %s", err)
	}
	return true, err
}
