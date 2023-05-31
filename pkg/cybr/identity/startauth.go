package identity

import (
	"fmt"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/identity/requests"
)

// StartAuthentication starts the authentication process
func StartAuthentication(c *api.Client, req requests.StartAuthentication) error {
	if c != nil {
		err := c.IsValid()
		if err != nil {
			return err
		}

		// Handle start of Identity authentication
		url := fmt.Sprintf("%s/Security/StartAuthentication", c.BaseURL)
		fmt.Printf(url)
		// Need to implement service specification in httpjson before continuing...
		if err != nil {
			return fmt.Errorf("Failed to start authentication. %s", err)
		}
	}

	return nil
}
