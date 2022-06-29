package ecs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators/aws"
)

// ECS aws resource type
type ECS struct {
	AwsCredentialRelativeURI string
	MetadataIP               string
}

func (r ECS) getCredentialURL() (string, error) {
	relativeURI := os.Getenv(r.AwsCredentialRelativeURI)
	if relativeURI == "" {
		return "", fmt.Errorf("Failed to get relative uri from environment variable '%s'", r.AwsCredentialRelativeURI)
	}
	return fmt.Sprintf("http://%s%s", r.MetadataIP, relativeURI), nil
}

// Name of aws resource
func (r ECS) Name() string {
	return "ecs"
}

// GetCredential from a ECS container
func (r ECS) GetCredential() (aws.Credential, error) {
	url, err := r.getCredentialURL()
	if err != nil {
		return aws.Credential{}, err
	}
	resp, err := http.Get(url)
	if err != nil {
		return aws.Credential{}, fmt.Errorf("Failed to retrieve AWS Credentials. %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return aws.Credential{}, fmt.Errorf("Failed to read AWS Credentials. %s", err)
	}

	var cred aws.Credential
	err = json.Unmarshal(body, &cred)
	if err != nil {
		return aws.Credential{}, fmt.Errorf("Failed to unmarshal IAM credential for URL '%s', %s", url, err)
	}

	return cred, nil
}

// New ECS resource
func New() ECS {
	return ECS{
		MetadataIP:               "169.254.170.2",
		AwsCredentialRelativeURI: "AWS_CONTAINER_CREDENTIALS_RELATIVE_URI",
	}
}
