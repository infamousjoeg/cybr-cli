package pasapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/infamousjoeg/pas-api-go/pasapi/httpJson"
)

const (
	// BaseURL is a constant exported variable to serve as the default
	// base URL when one is not provided.
	BaseURL = "https://localhost/PasswordVault"
)

func header(token string) http.Header {
	header := make(http.Header)
	header.Add("Authorization", token)
	return header
}

// Authenticate to the PAS API which requires you to provide
// Base URL (https://pvwa.example.com/PasswordVault), username, password,
// and authentication type (cyberark, ldap, radius, saml).
func Authenticate(baseURL string, username string, password string, authType string) (token string, err error) {
	// Set baseURL to default constant if not provided
	if baseURL == "" {
		baseURL = BaseURL
	}
	// Format the URI for the request
	url := fmt.Sprintf("https://%s/API/auth/%s/Logon", baseURL, authType)
	// Format the body for the request
	body := fmt.Sprintf("{\"username\": \"%s\", \"password\":\"%s\"}", username, password)
	// POST the request in raw form
	response, err := httpJson.SendRequestRaw(url, "POST", nil, body)
	if err != nil {
		return "", fmt.Errorf("Failed to authenticate to the PAS API. %s", err)
	}
	return strings.Trim(string(response), "\""), err
}

// ListApplications from the PAS API where all Application IDs are returned.
func ListApplications(hostname string, token string, location string, includeSublocations bool) (map[string]interface{}, error) {
	params := url.Values{}
	// If location is provided, set it - otherwise, default is `\`
	if location != "" {
		params.Add("Location",location)
	}
	// If includeSublocations is provided, set it - otherwise, default is `true`
	if includeSublocations != nil or includeSublocations != true {
		params.Add("IncludeSublocations","false")
	}
	// Format the URI for the request and provide URL parameters
	url := fmt.Sprintf("https://%s/WebServices/PIMServices.svc/Applications", hostname, params)
	// Format the header properly
	header := header(token)
	// GET the request
	response, err := httpJson.Get(url, header)
	return response, err
}

/*
// ListApplicationAuthenticationMethods from cyberark
func ListApplicationAuthenticationMethods(hostname string, token string, appName string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://%s/PasswordVault/WebServices/PIMServices.svc/Applications/%s/Authentications", hostname, appName)
	header := header(token)
	response, err := httpJson.Get(url, header)
	return response, err
}

// ListSafes from cyberark user has access too
func ListSafes(hostname string, token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://%s/PasswordVault/WebServices/PIMServices.svc/Safes", hostname)
	// fmt.Println(url)
	header := make(http.Header)
	header.Add("Authorization", token)
	response, err := httpJson.Get(url, header)
	return response, err
}

// ListSafeMembers List all members of a safe
func ListSafeMembers(hostname string, token string, safeName string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://%s/PasswordVault/WebServices/PIMServices.svc/Safes/%s/Members", hostname, safeName)
	// fmt.Println(url)
	header := make(http.Header)
	header.Add("Authorization", token)
	response, err := httpJson.Get(url, header)
	return response, err
}

// GetSafesUserIsMemberOf Iterate through all safes and see which safe username is member of
// This method performs ListSafes() and ListSafeMembers()
func GetSafesUserIsMemberOf(hostname string, token string, username string) ([]string, error) {
	// List safes
	response, err := ListSafes(hostname, token)
	if err != nil {
		return nil, fmt.Errorf("Failed to list cyberark safes %s", err)
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
*/