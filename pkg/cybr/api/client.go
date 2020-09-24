package api

import (
	"fmt"
)

// Client contains the data necessary for requests to pass successfully
type Client struct {
	Hostname     string
	AuthType     AuthType
	sessionToken string
}

// AuthType is the authentication method that will be used for Logon
type AuthType string

const (
	// CyberArk is referencing the case switch
	CyberArk AuthType = "cyberark"
	// LDAP is referencing the case switch
	LDAP = "ldap"
)

// IsValid checks to make sure that the authentication method chosen is valid
func (lt AuthType) IsValid() error {
	switch lt {
	case CyberArk, LDAP:
		return nil
	}
	return fmt.Errorf("Invalid auth type '%s'", lt)
}
