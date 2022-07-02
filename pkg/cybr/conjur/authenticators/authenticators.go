package authenticators

import (
	"fmt"
	"strings"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/conjur/authenticators/iam"
	helpersauthn "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators"
)

func getAuthenticators(config helpersauthn.Config) []helpersauthn.Authenticator {
	authenticators := []helpersauthn.Authenticator{}
	authenticators = append(authenticators, iam.New())
	return authenticators
}

// GetAuthenticator will return the authenticator client for the given name
func GetAuthenticator(name string, config helpersauthn.Config) (helpersauthn.Authenticator, error) {
	authenticators := getAuthenticators(config)
	for _, r := range authenticators {
		if strings.ToLower(name) == r.Name() {
			return r, nil
		}
	}

	return nil, fmt.Errorf("Failed to retrieve authenticator with name '%s'", name)
}
