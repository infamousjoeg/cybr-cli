package api

import (
	"fmt"
	"strings"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// Logon to PAS REST API Web Service
// Because we're using concurrentSession capability, this is only supported
// on PAS REST API v11.3 and above
func (c *Client) Logon(req requests.Logon) error {
	err := c.IsValid()
	if err != nil {
		return err
	}

	// Handle cyberark, ldap, and radius push, append & challenge/response authentication methods
	url := fmt.Sprintf("%s/passwordvault/api/auth/%s/logon", c.BaseURL, c.AuthType)
	token, err := httpJson.SendRequestRaw(url, "POST", "", req, c.InsecureTLS, c.Logger)
	if err != nil {
		return fmt.Errorf("Failed to authenticate to the PAS REST API. %s", err)
	}

	c.SessionToken = strings.Trim(string(token), "\"")
	return nil
}

// Logoff the PAS REST API Web Service
func (c Client) Logoff() error {
	// Set URL for request
	url := fmt.Sprintf("%s/passwordvault/api/auth/logoff", c.BaseURL)
	_, err := httpJson.Post(url, c.SessionToken, nil, c.InsecureTLS, c.Logger)
	if err != nil {
		return fmt.Errorf("Unable to logoff PAS REST API Web Service. %s", err)
	}
	return nil
}
