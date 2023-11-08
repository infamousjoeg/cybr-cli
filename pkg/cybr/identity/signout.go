package identity

import (
	"fmt"
	"net/http"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// SignOutSession signs out of the current Identity session
func SignOutSession(c api.Client) error {
	identityTenant := fmt.Sprintf("https://%s.id.cyberark.cloud", c.TenantID)
	url := fmt.Sprintf("%s/UserMgmt/SignOutCurrentSession", identityTenant)

	headers := http.Header{}
	headers.Add("X-IDAP-NATIVE-CLIENT", "true")
	headers.Add("Content-Type", "application/json")

	_, err := httpjson.SendRequestRawWithHeaders(url, "POST", headers, nil, c.InsecureTLS, c.Logger)
	if err != nil {
		return fmt.Errorf("Failed to sign out of current session. %s", err)
	}

	return nil
}
