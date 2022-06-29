package iam

import (
	"github.com/cyberark/conjur-api-go/conjurapi"
)

type IAM struct {
	GetAuthnUrl      string
	GetAuthnResource string
}

// Name of the authenticator type
func (r IAM) Name() string {
	return "authn-iam"
}

// GetConfig will retrieve the details of the authenticator
func (r IAM) Authenticate() (*conjurapi.Client, error) {
	return &conjurapi.Client{}, nil
}

// Authenticate will create a new authn-iam Conjur Client
func New() IAM {
	return IAM{
		GetAuthnUrl:      "http://169.254.169.254/latest/meta-data/iam/security-credentials/",
		GetAuthnResource: "http://169.254.169.254/latest/meta-data/iam/security-credentials/%s",
	}
}
