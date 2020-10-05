package api

import (
	"encoding/json"
	"fmt"

	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// VerifyResponse contains all response data from /verify endpoint
type VerifyResponse struct {
	ApplicationName       string                  `json:"ApplicationName"`
	AuthenticationMethods []AuthenticationMethods `json:"AuthenticationMethods"`
	ServerID              string                  `json:"ServerId"`
	ServerName            string                  `json:"ServerName"`
}

// AuthenticationMethods feeds into VerifyResponse
type AuthenticationMethods struct {
	Enabled bool   `json:"Enabled"`
	ID      string `json:"Id"`
}

// ServerVerify is an unauthenticated endpoint for testing Web Service availability
func (c Client) ServerVerify() (*VerifyResponse, error) {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Verify", c.BaseURL)
	response, err := httpJson.SendRequest(url, "GET", "", nil)
	if err != nil {
		return &VerifyResponse{}, fmt.Errorf("Error verifying PAS REST API Web Service. %s", err)
	}
	jsonString, _ := json.Marshal(response)
	VerifyResponse := VerifyResponse{}
	err = json.Unmarshal(jsonString, &VerifyResponse)
	return &VerifyResponse, err
}
