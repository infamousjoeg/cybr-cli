package ec2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators/aws"
)

// EC2 holds the relevant internal metadata URLs for AWS EC2
type EC2 struct {
	GetIamRoleURL       string
	GetIamCredentialURL string
}

func (r EC2) getIamRoleName() (string, error) {
	resp, err := http.Get(r.GetIamRoleURL)
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve IAM role name. %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read IAM role name. %s", err)
	}

	return string(body), nil
}

func (r EC2) getIamCredential(iamRoleName string) (aws.Credential, error) {
	resp, err := http.Get(fmt.Sprintf(r.GetIamCredentialURL, iamRoleName))
	if err != nil {
		return aws.Credential{}, fmt.Errorf("Failed to retrieve AM credential for role '%s'. %s", iamRoleName, err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return aws.Credential{}, fmt.Errorf("Failed to read IAM credential for role '%s'. %s", iamRoleName, err)
	}

	var cred aws.Credential
	err = json.Unmarshal(body, &cred)
	if err != nil {
		return aws.Credential{}, fmt.Errorf("Failed to unmarshal IAM credential for role '%s', %s", iamRoleName, err)
	}

	return cred, nil
}

// Name of the resource type
func (r EC2) Name() string {
	return "ec2"
}

// GetCredential will retrieve an IAM credential
func (r EC2) GetCredential() (aws.Credential, error) {
	iamRoleName, err := r.getIamRoleName()
	if err != nil {
		return aws.Credential{}, err
	}
	return r.getIamCredential(iamRoleName)
}

// New will create a new EC2 AWS Resource
func New() EC2 {
	return EC2{
		GetIamRoleURL:       "http://169.254.169.254/latest/meta-data/iam/security-credentials/",
		GetIamCredentialURL: "http://169.254.169.254/latest/meta-data/iam/security-credentials/%s",
	}
}
