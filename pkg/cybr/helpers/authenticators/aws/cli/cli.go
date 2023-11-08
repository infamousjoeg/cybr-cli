package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators/aws"
)

var iamSessionName = "cybr-cli"

// CLI holds the relevant internal metadata URLs for AWS EC2
type CLI struct {
	GetIamRoleURL       string
	GetIamCredentialURL string
}

func (r CLI) getIamRoleName() (string, error) {
	resp := os.Getenv("CONJUR_AUTHN_HOST_ID")

	splitResp := strings.Split(resp, "/")
	// Extract AWS Account ID from the value
	awsAccountID := splitResp[len(splitResp)-2]
	// Extract IAM Role Name from the value
	iamRoleName := splitResp[len(splitResp)-1]
	// Create the IAM Role ARN
	iamRoleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", awsAccountID, iamRoleName)

	return iamRoleArn, nil
}

func (r CLI) getIamCredential(iamRoleName string) (aws.Credential, error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return aws.Credential{}, fmt.Errorf("Failed to retrieve IAM credential for role '%s'. %s", iamRoleName, err)
	}

	// Create an AWS STS client from the config
	client := sts.NewFromConfig(cfg)

	// Assume the role
	input := &sts.AssumeRoleInput{
		RoleArn:         &iamRoleName,
		RoleSessionName: &iamSessionName,
	}
	result, err := client.AssumeRole(context.TODO(), input)

	var cred aws.Credential
	cred.AccessKeyID = *result.Credentials.AccessKeyId
	cred.SecretAccessKey = *result.Credentials.SecretAccessKey
	cred.Token = *result.Credentials.SessionToken

	return cred, nil
}

// Name of the resource type
func (r CLI) Name() string {
	return "cli"
}

// GetCredential will retrieve an IAM credential
func (r CLI) GetCredential() (aws.Credential, error) {
	iamRoleName, err := r.getIamRoleName()
	if err != nil {
		return aws.Credential{}, err
	}
	return r.getIamCredential(iamRoleName)
}

// New will create a new EC2 AWS Resource
func New() CLI {
	return CLI{}
}
