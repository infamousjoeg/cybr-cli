package authenticators

import (
	"github.com/cyberark/conjur-api-go/conjurapi"
)

// Config holds the configuration for the Conjur authenticator
type Config struct {
	Account         string
	ApplianceURL    string
	Login           string
	ServiceID       string
	IgnoreSSLVerify bool
}

// Authenticator is used to retrieve the Conjur client.
// This is required because authn uses username and password to authenticate to Conjur, while authn-iam uses a token.
type Authenticator interface {
	Name() string
	Authenticate(config Config) (*conjurapi.Client, error)
}

// GetAuthURL returns a proper LDAP Authentication authn_url for the ~/.conjurrc file
func GetAuthURL(baseURL string, authType string, serviceID string) string {
	authURL := baseURL
	if authType != "" {
		authURL = authURL + "/" + authType
	}
	if serviceID != "" {
		authURL = authURL + "/" + serviceID
	}
	return authURL
}
