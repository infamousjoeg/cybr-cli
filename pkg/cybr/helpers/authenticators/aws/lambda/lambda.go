package lambda

import (
	"fmt"
	"os"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators/aws"
)

// Lambda represents a lambda function and the environment variables where the AWS credentials are located
type Lambda struct {
	AccessKey       string
	SecretAccessKey string
	Token           string
}

// Name of the AWS resource
func (r Lambda) Name() string {
	return "lambda"
}

// GetCredential get the IAM credentials from the lambda environment variables
func (r Lambda) GetCredential() (aws.Credential, error) {
	accessKey := os.Getenv(r.AccessKey)
	if accessKey == "" {
		return aws.Credential{}, fmt.Errorf("Failed to retrieve access key from environment variables '%s'", r.AccessKey)
	}

	secretAccessKey := os.Getenv(r.SecretAccessKey)
	if secretAccessKey == "" {
		return aws.Credential{}, fmt.Errorf("Failed to retrieve secret access key from environment variables '%s'", r.SecretAccessKey)

	}

	token := os.Getenv(r.Token)
	if token == "" {
		return aws.Credential{}, fmt.Errorf("Failed to retrieve token from environment variables '%s'", r.Token)
	}

	return aws.Credential{
		AccessKeyID:     accessKey,
		SecretAccessKey: secretAccessKey,
		Token:           token,
	}, nil
}

// New returns a new Lambda object
func New() Lambda {
	return Lambda{
		AccessKey:       "AWS_ACCESS_KEY_ID",
		SecretAccessKey: "AWS_SECRET_ACCESS_KEY",
		Token:           "AWS_SESSION_TOKEN",
	}
}
