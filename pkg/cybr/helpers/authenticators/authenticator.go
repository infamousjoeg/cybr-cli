package authenticators

import (
	"github.com/cyberark/conjur-api-go/conjurapi"
)

// Authenticator is used to retrieve the Conjur client.
// This is required because authn uses username and password to authenticate to Conjur, while authn-iam uses a token.
type Authenticator interface {
	Name() string
	Authenticate() (*conjurapi.Client, error)
}
