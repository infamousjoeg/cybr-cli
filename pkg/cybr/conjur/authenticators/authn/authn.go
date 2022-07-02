package authn

import "github.com/cyberark/conjur-api-go/conjurapi"

// Authn is the struct that holds the details of the authentication
type Authn struct {
	GetAuthnURL      string
	GetAuthnResource string
}

// Name of the authenticator type
func (r Authn) Name() string {
	return "authn"
}

// Authenticate will authenticate the user and return a Conjur API client
func (r Authn) Authenticate() (*conjurapi.Client, error) {
	return &conjurapi.Client{}, nil
}

// New will returns a new authn object
func New() Authn {
	return Authn{
		GetAuthnURL:      "http://169.254.169.254/latest/meta-data/iam/security-credentials/",
		GetAuthnResource: "http://169.254.169.254/latest/meta-data/iam/security-credentials/%s",
	}
}
