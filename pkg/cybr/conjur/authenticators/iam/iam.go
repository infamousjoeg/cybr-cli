package iam

import (
	"fmt"
	"os"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators"
	helpersauthn "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators/aws"
)

// IAM represents the config for IAM authenticator
type IAM struct {
	Account      string
	ApplianceURL string
	Login        string
	ServiceID    string
}

// Name of the authenticator type
func (r IAM) Name() string {
	return "authn-iam"
}

// Authenticate will create a new authn-iam Conjur Client
func (r IAM) Authenticate(config helpersauthn.Config) (*conjurapi.Client, error) {
	// Get the AWS Service Type (EC2, ECS, or Lambda)
	awsServiceType := os.Getenv("CONJUR_AWS_TYPE")
	if awsServiceType != "EC2" && awsServiceType != "ECS" && awsServiceType != "Lambda" {
		return nil, fmt.Errorf("CONJUR_AWS_TYPE environment variable is not set or is not set to EC2, ECS, or Lambda")
	}

	// Get metadata URLs for the AWS Service
	resource, err := GetAwsResource(awsServiceType)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve AWS resource with type '%s'. %s", awsServiceType, err)
	}

	// Get the IAM credentials from AWS STS
	credential, err := resource.GetCredential()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve AWS credential for service type '%s'. %s", awsServiceType, err)
	}

	// Get the Signed Headers for the IAM credentials from AWS STS
	conjurAuthnRequest, err := aws.GetAuthenticationRequestNow(credential.AccessKeyID, credential.SecretAccessKey, credential.Token)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve AWS authentication request for service type '%s'. %s", awsServiceType, err)
	}

	// Create the URL for the authenticator
	authnURL := authenticators.GetAuthURL(config.ApplianceURL, "authn-iam", config.ServiceID)

	// Authenticate to Conjur using the AWS STS signed headers and receive a session token
	accessToken, err := Authenticate(authnURL, config.Account, config.Login, conjurAuthnRequest, false, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to authenticate to Conjur with service type '%s'. %s", awsServiceType, err)
	}

	// Create the necessary config for the Conjur Client
	conjurConfig := conjurapi.Config{
		Account:      config.Account,
		ApplianceURL: config.ApplianceURL,
	}

	// Create the Conjur Client using the Conjur API session token
	client, err := conjurapi.NewClientFromToken(conjurConfig, string(accessToken))
	if err != nil {
		return nil, fmt.Errorf("Failed to create Conjur Client with service type '%s'. %s", awsServiceType, err)
	}

	return client, nil
}

// New returns a new IAM object
func New() IAM {
	return IAM{
		Account:      "",
		ApplianceURL: "",
		Login:        "",
		ServiceID:    "",
	}
}
