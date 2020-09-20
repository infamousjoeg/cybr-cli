package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	httpJson "github.com/infamousjoeg/pas-api-go/pkg/helpers"
)

func header(token string) http.Header {
	header := make(http.Header)
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", token)
	return header
}

// Authenticate to PAS REST API /logon endpoint
// Because we're using concurrentSession capability, this is only supported
// on PAS REST API v11.3 and above
func Authenticate(hostname string, username string, password string, authType string) (token string, err error) {
	// Return error if unsupported authentication type chosen
	if authType != "cyberark" && authType != "ldap" {
		log.Fatal("Unsupported auth type used. Only 'cyberark' or 'ldap' supported")
	}

	// Set URL for request
	url := fmt.Sprintf("%s/PasswordVault/api/auth/%s/logon", hostname, authType)
	body := fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\"}", username, password)

	// Send request and received response
	response, err := httpJson.SendRequestRaw(url, "POST", nil, body)
	if err != nil {
		return "", fmt.Errorf("Failed to authenticate to the PAS REST API. %s", err)
	}
	// Trim "..." from response and return session token
	return strings.Trim(string(response), "\""), err
}

// ListApplications from cyberark
func ListApplications(hostname string, token string, location string, subLocations bool) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications?Location=%s&IncludeSublocations=%t", hostname, location, subLocations)
	header := header(token)
	response, err := httpJson.Get(url, header)
	return response, err
}

// ListApplicationAuthenticationMethods from cyberark
func ListApplicationAuthenticationMethods(hostname string, token string, appID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications/%s/Authentications", hostname, appID)
	header := header(token)
	response, err := httpJson.Get(url, header)
	return response, err
}

// ListSafes from cyberark user has access too
func ListSafes(hostname string, token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/PasswordVault/api/safes", hostname)
	header := header(token)
	response, err := httpJson.Get(url, header)
	return response, err
}

// ListSafeMembers List all members of a safe
func ListSafeMembers(hostname string, token string, safeName string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Safes/%s/Members", hostname, safeName)
	header := header(token)
	response, err := httpJson.Get(url, header)
	return response, err
}

// GetSafesUserIsMemberOf Iterate through all safes and see which safe username is member of
// This method performs ListSafes() and ListSafeMembers()
func GetSafesUserIsMemberOf(hostname string, token string, username string) ([]string, error) {
	// List safes
	response, err := ListSafes(hostname, token)
	if err != nil {
		return nil, fmt.Errorf("Failed to list CyberArk safes %s", err)
	}

	safes := response["GetSafesResult"].([]interface{})
	var safesUserIsMemberOf []string

	// Iterate through safes get members and see if member of safe is 'username' provided
	for _, safe := range safes {
		safeInterface := safe.(map[string]interface{})
		safeName := safeInterface["SafeName"].(string)

		// now query each safe for this specific username
		response, err = ListSafeMembers(hostname, token, safeName)
		if err != nil {
			return nil, fmt.Errorf("Failed to list safe members for safe '%s'. %s", safeName, err)
		}

		members := response["members"].([]interface{})

		for _, member := range members {
			memberInterface := member.(map[string]interface{})
			safeMember := memberInterface["UserName"].(string)
			if safeMember == username {
				safesUserIsMemberOf = append(safesUserIsMemberOf, safeName)
			}
		}
	}
	return safesUserIsMemberOf, err
}

// ServerVerify is an unauthenticated endpoint for testing Web Service availability
func ServerVerify(hostname string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Verify", hostname)
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	response, err := httpJson.Get(url, header)
	return response, err
}
