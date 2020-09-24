package api

import (
	"fmt"
)

type Client struct {
	Hostname     string
	AuthType     AuthType
	sessionToken string
}

type AuthType string

const (
	Cyberark AuthType = "cyberark"
	LDAP              = "ldap"
)

func (lt AuthType) IsValid() error {
	switch lt {
	case Cyberark, LDAP:
		return nil
	}
	return fmt.Errorf("Invalid auth type '%s'", lt)
}
