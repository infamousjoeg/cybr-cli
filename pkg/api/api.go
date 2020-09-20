package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	httpJson "github.com/infamousjoeg/pas-api-go/pkg/helpers"
)

func header(token string) http.Header {
	header := make(http.Header)
	header.Add("Authorization", token)
	return header
}

// AuthParams contains parameters for Authenticate
type AuthParams struct {
	Hostname string
	AuthType string
}

// AuthRequest contains body data for Authenticate
type AuthRequest struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	NewPassword        string `json:"newPassword,omitempty"`
	ConcurrentSessions bool   `json:"concurrentSessions,omitempty"`
}

// Authenticate to PAS REST API /logon endpoint
func Authenticate(params *AuthParams, data *AuthRequest) (token string, err error) {
	// Convert hostname to lowercase or else set to default 'components.cyberarkdemo.com'
	var hostname string
	if params.Hostname != "" {
		hostname = strings.ToLower(params.Hostname)
	} else {
		hostname = "https://components.cyberarkdemo.com"
	}
	// Convert authType to lowercase or else set to default 'cyberark'
	var authType string
	if params.AuthType != "" {
		authType = strings.ToLower(params.AuthType)
	} else {
		authType = "cyberark"
	}

	// Return error if unsupported authentication type chosen
	if authType != "cyberark" && authType != "ldap" {
		return "", fmt.Errorf("Unsupported auth type. Only 'cyberark' or 'ldap' supported")
	}

	url := fmt.Sprintf("%s/PasswordVault/api/auth/%s/logon", hostname, authType)

	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	response, err := httpJson.SendRequestRaw(url, "POST", nil, string(body))
	if err != nil {
		return "", fmt.Errorf("Failed to authenticate to the cyberark api. %s", err)
	}
	return strings.Trim(string(response), "\""), err
}

// ListApplications from cyberark
func ListApplications(hostname string, token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://%s/PasswordVault/WebServices/PIMServices.svc/Applications/", hostname)
	header := header(token)
	response, err := httpJson.Get(url, header)
	return response, err
}

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
