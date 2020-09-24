package api

import (
	"fmt"
)

// Client contains the data necessary for requests to pass successfully
type Client struct {
	BaseURL      string
	AuthType     string
	sessionToken string
}

// IsValid checks to make sure that the authentication method chosen is valid
func (c *Client) IsValid() error {
	if c.AuthType == "cyberark" || c.AuthType == "ldap" {
		return nil
	}
	return fmt.Errorf("Invalid auth type '%s'", c.AuthType)
}
