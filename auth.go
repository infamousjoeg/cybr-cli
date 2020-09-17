package pasapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SessionToken contains the session token value to return
type SessionToken struct {
	Token string
}

// AuthOptions contains all the variables needed to complete a logon or logoff request
type AuthOptions struct {
	Action            string
	AuthType          string
	APIUsername       string
	APIPassword       []byte
	NewPassword       []byte
	ConcurrentSession bool
}

// Authentication sends a POST request to logon or logoff the PAS REST API
func (c *Client) Authentication(options *AuthOptions) (*SessionToken, error) {
	// Declare variables
	var action string
	var authType string
	var apiUsername string
	var apiPassword []byte
	var newPassword []byte
	var concurrentSession bool

	// Add options to declared variables if not nil
	if options != nil {
		action = strings.ToLower(options.Action)
		authType = strings.ToLower(options.AuthType)
		apiUsername = options.APIUsername
		apiPassword = options.APIPassword
		newPassword = options.NewPassword
		concurrentSession = options.ConcurrentSession
	}

	if action != "logon" || action != "logoff" {
		fmt.Errorf("It is required that you declare an action (logon or logoff)")
	}

	if authType != "cyberark" || authType != "ldap" {
		fmt.Errorf("You have not declared a valid authType (cyberark or ldap)")
	}

	// Begin building request body
	body := url.Values{}
	if apiUsername != "" && apiPassword != nil {
		body.Set("username", apiUsername)
		body.Set("password", string(apiPassword))
	} else {
		fmt.Errorf("You must provide a valid API Username and Password")
	}
	if newPassword != nil {
		body.Set("newPassword", string(newPassword))
	}
	if concurrentSession != false {
		body.Set("concurrentSession", "true")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/auth/%s/%s", c.BaseURL, authType, action), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json charset=utf-8")

	res := SessionToken{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
