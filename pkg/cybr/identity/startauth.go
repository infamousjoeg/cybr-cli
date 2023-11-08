package identity

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/identity/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/identity/responses"
)

var identityTenant string

// StartAuthentication starts the authentication process
func StartAuthentication(c api.Client, req requests.StartAuthentication, podFqdn string) (*responses.Authentication, error) {
	if podFqdn == "" {
		identityTenant = fmt.Sprintf("https://%s.id.cyberark.cloud", req.TenantID)
	} else {
		identityTenant = fmt.Sprintf("https://%s", podFqdn)
	}
	url := fmt.Sprintf("%s/Security/StartAuthentication", identityTenant)

	headers := http.Header{}
	headers.Add("X-IDAP-NATIVE-CLIENT", "true")
	headers.Add("Content-Type", "application/json")

	res, err := httpjson.SendRequestRawWithHeaders(url, "POST", headers, req, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.Authentication{}, fmt.Errorf("Failed to start authentication. %s", err)
	}

	StartAuthResponse := &responses.Authentication{}
	err = json.Unmarshal(res, StartAuthResponse)
	return StartAuthResponse, err
}
