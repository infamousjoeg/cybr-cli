package authenticators

import (
	"fmt"
	"strings"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/conjur/authenticators/authn"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/conjur/authenticators/iam"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators"
)

func getAuthenticators() []authenticators.Authenticator {
	authenticators := []authenticators.Authenticator{}
	authenticators = append(authenticators, iam.New())
	authenticators = append(authenticators, authn.New())
	return authenticators
}

// GetAuthenticator will return the authenticator client for the given name
func GetAuthenticator(name string) (authenticators.Authenticator, error) {
	authenticators := getAuthenticators()
	for _, r := range authenticators {
		if strings.ToLower(name) == r.Name() {
			return r, nil
		}
	}

	return nil, fmt.Errorf("Failed to retrieve authenticator with name '%s'", name)
}
