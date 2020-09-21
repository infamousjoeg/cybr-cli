package api

import (
	"fmt"

	httpJson "github.com/infamousjoeg/pas-api-go/pkg/helpers"
)

// ListApplications from cyberark
func ListApplications(hostname string, token string, location string, subLocations bool) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications?Location=%s&IncludeSublocations=%t", hostname, location, subLocations)
	response, err := httpJson.Get(url, token)
	return response, err
}

// ListApplicationAuthenticationMethods from cyberark
func ListApplicationAuthenticationMethods(hostname string, token string, appID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications/%s/Authentications", hostname, appID)
	response, err := httpJson.Get(url, token)
	return response, err
}

// ListSafes from cyberark user has access too
func ListSafes(hostname string, token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/PasswordVault/api/safes", hostname)
	response, err := httpJson.Get(url, token)
	return response, err
}

// ListSafeMembers List all members of a safe
func ListSafeMembers(hostname string, token string, safeName string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Safes/%s/Members", hostname, safeName)
	response, err := httpJson.Get(url, token)
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
