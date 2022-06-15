package iam

import (
	"fmt"
	"strings"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/aws"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/aws/ec2"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/aws/ecs"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/aws/lambda"
)

func getAwsResources() []aws.AwsResource {
	resources := []aws.AwsResource{}
	resources = append(resources, ec2.New())
	resources = append(resources, lambda.New())
	resources = append(resources, ecs.New())
	return resources
}

// GetAwsResource will return an interface that has the ability to retrieve IAM AWS credentials from the desired metadata endpoint
func GetAwsResource(name string) (aws.AwsResource, error) {
	resources := getAwsResources()
	for _, r := range resources {
		if strings.ToLower(name) == r.Name() {
			return r, nil
		}
	}

	return nil, fmt.Errorf("Failed to retrieve AWS resource with type '%s'", name)
}
