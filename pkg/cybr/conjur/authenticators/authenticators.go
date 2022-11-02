package authenticators

import (
	"fmt"
	"strings"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/conjur/authenticators/iam"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators"
)

func getAuthenticators(config authenticators.Config) ([]authenticators.Authenticator, error) {
	authenticators := []authenticators.Authenticator{}
	iamInterface, err := iam.New()
	if err != nil {
		return nil, fmt.Errorf("Failed to create IAM authenticator. %s", err)
	}
	authenticators = append(authenticators, iamInterface)
	return authenticators, nil
}

// GetAuthenticator will return the authenticator client for the given name
func GetAuthenticator(name string, config authenticators.Config) (authenticators.Authenticator, error) {
	authenticators, err := getAuthenticators(config)
	if err != nil {
		return nil, err
	}
	for _, r := range authenticators {
		if strings.ToLower(name) == r.Name() {
			return r, nil
		}
	}

	return nil, fmt.Errorf("Failed to retrieve authenticator with name '%s'", name)
}
