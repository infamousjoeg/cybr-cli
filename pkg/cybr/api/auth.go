package api

import (
	"fmt"
	"strings"

	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

var (
	token []byte
)

// LogonRequest contains the body of the Logon function's request
type LogonRequest struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	ConcurrentSession bool   `json:"concurrentSession,omitempty"`
}

// Logon to PAS REST API Web Service
// Because we're using concurrentSession capability, this is only supported
// on PAS REST API v11.3 and above
func (c *Client) Logon(req LogonRequest) error {
	err := c.IsValid()
	if err != nil {
		return err
	}

	// Handle cyberark, ldap, and radius push & append methods authentication methods
	url := fmt.Sprintf("%s/PasswordVault/api/auth/%s/logon", c.BaseURL, c.AuthType)
	token, err = httpJson.SendRequestRaw(url, "POST", "", req, c.InsecureTLS)
	if err != nil {
		return fmt.Errorf("Failed to authenticate to the PAS REST API. %s", err)
	}

	c.SessionToken = strings.Trim(string(token), "\"")
	return nil
}

// Logoff the PAS REST API Web Service
func (c Client) Logoff() error {
	// Set URL for request
	url := fmt.Sprintf("%s/PasswordVault/api/auth/logoff", c.BaseURL)
	_, err := httpJson.Post(url, c.SessionToken, nil, c.InsecureTLS)
	if err != nil {
		return fmt.Errorf("Unable to logoff PAS REST API Web Service. %s", err)
	}
	return nil
}
