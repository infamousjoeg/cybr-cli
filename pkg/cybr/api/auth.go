package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/shared"
	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// Logon to PAS REST API Web Service
// Because we're using concurrentSession capability, this is only supported
// on PAS REST API v11.3 and above
func (c *Client) Logon(ctx context.Context, req requests.Logon) (context.Context, *shared.ErrorResponse, error) {
	err := c.IsValid()
	if err != nil {
		return ctx, &shared.ErrorResponse{}, err
	}

	// Handle cyberark, ldap, and radius push, append & challenge/response authentication methods
	url := fmt.Sprintf("%s/passwordvault/api/auth/%s/logon", c.BaseURL, c.AuthType)
	ctx, response, err := httpJson.SendRequestRaw(ctx, false, url, "POST", "", req, c.InsecureTLS, c.Logger)
	if err != nil || strings.Contains(string(response), "ITATS542I") {
		// Check if response can be unmarshalled to ErrorResponse
		errorResponse := &shared.ErrorResponse{}
		errUm := json.Unmarshal(response, &errorResponse)
		if errUm != nil {
			return ctx, nil, errUm
		}

		return ctx, errorResponse, fmt.Errorf("Failed to authenticate to the PAS REST API. %s", err)
	}

	c.SessionToken = strings.Trim(string(response), "\"")
	return ctx, nil, nil
}

// Logoff the PAS REST API Web Service
func (c Client) Logoff() error {
	// Set URL for request
	url := fmt.Sprintf("%s/passwordvault/api/auth/logoff", c.BaseURL)
	_, err := httpJson.Post(false, url, c.SessionToken, nil, c.InsecureTLS, c.Logger)
	if err != nil {
		return fmt.Errorf("Unable to logoff PAS REST API Web Service. %s", err)
	}
	return nil
}
