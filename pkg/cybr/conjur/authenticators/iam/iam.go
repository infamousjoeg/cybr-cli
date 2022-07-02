package iam

import (
	"github.com/cyberark/conjur-api-go/conjurapi"
)

// IAM represents the IAM authenticator
type IAM struct {
	GetAuthnURL      string
	GetAuthnResource string
}

// Name of the authenticator type
func (r IAM) Name() string {
	return "authn-iam"
}

// Authenticate will create a new authn-iam Conjur Client
func (r IAM) Authenticate() (*conjurapi.Client, error) {
	return &conjurapi.Client{}, nil
}

// New returns a new IAM object
func New() IAM {
	return IAM{
		GetAuthnURL:      "http://169.254.169.254/latest/meta-data/iam/security-credentials/",
		GetAuthnResource: "http://169.254.169.254/latest/meta-data/iam/security-credentials/%s",
	}
}
