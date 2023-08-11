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

// AdvanceAuthentication will answer challenges from CyberArk Identity
func AdvanceAuthentication(c api.Client, req requests.AdvanceAuthentication) (*responses.Authentication, error) {
	identityTenant := fmt.Sprintf("https://%s.id.cyberark.cloud", c.TenantID)
	url := fmt.Sprintf("%s/Security/AdvanceAuthentication", identityTenant)

	headers := http.Header{}
	headers.Add("X-IDAP-NATIVE-CLIENT", "true")
	headers.Add("Content-Type", "application/json")

	res, err := httpjson.SendRequestRawWithHeaders(url, "POST", headers, req, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.Authentication{}, fmt.Errorf("Failed to start authentication. %s", err)
	}

	AdvanceAuthResponse := &responses.Authentication{}
	err = json.Unmarshal(res, AdvanceAuthResponse)
	return AdvanceAuthResponse, err
}
