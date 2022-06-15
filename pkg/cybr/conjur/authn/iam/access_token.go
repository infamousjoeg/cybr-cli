package iam

import (
	"io/ioutil"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/aws"
)

// GetConjurAccessToken Get Conjur access token from Conjur
func GetConjurAccessToken(config Config) ([]byte, error) {
	resource, err := GetAwsResource(config.AWSName)
	if err != nil {
		return nil, err
	}

	// Retrieveing the AWS IAM Credential
	credential, err := resource.GetCredential()
	if err != nil {
		return nil, err
	}

	// Convert the AWS IAM Credential into a Conjur Authentication request
	conjurAuthnRequest, err := aws.GetAuthenticationRequestNow(credential.AccessKeyID, credential.SecretAccessKey, credential.Token)
	if err != nil {
		return nil, err
	}

	// Use the Authentication request to authenticate to Conjur and get a Conjur access token
	cert, err := config.getCertificate()
	if err != nil {
		return nil, err
	}
	accessToken, err := Authenticate(config.AuthnURL, config.Account, config.Login, conjurAuthnRequest, config.IgnoreSSLVerify, cert)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

// WriteAccessToken witll write Conjur access token to a file specified
func WriteAccessToken(accessToken []byte, tokenPath string) error {
	if tokenPath == "" {
		return nil
	}

	err := ioutil.WriteFile(tokenPath, accessToken, 0400)
	if err != nil {
		return err
	}

	return nil
}
