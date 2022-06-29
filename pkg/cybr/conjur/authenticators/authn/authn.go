package authn

import "github.com/cyberark/conjur-api-go/conjurapi"

type Authn struct {
	GetAuthnUrl      string
	GetAuthnResource string
}

// Name of the authenticator type
func (r Authn) Name() string {
	return "authn"
}

// GetConfig will retrieve the details of the authenticator
func (r Authn) Authenticate() (*conjurapi.Client, error) {
	return &conjurapi.Client{}, nil
}

// New will create a new authn Conjur Client
func New() Authn {
	return Authn{
		GetAuthnUrl:      "http://169.254.169.254/latest/meta-data/iam/security-credentials/",
		GetAuthnResource: "http://169.254.169.254/latest/meta-data/iam/security-credentials/%s",
	}
}
