package api

import (
	"encoding/json"
	"fmt"
	"log"

	httpJson "github.com/infamousjoeg/pas-api-go/pkg/cybr/helpers"
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
func ServerVerify(hostname string) *VerifyResponse {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Verify", hostname)
	response, err := httpJson.SendRequest(url, "GET", "", "")
	if err != nil {
		log.Fatalf("Error verifying PAS REST API Web Service. %s", err)
	}
	jsonString, _ := json.Marshal(response)
	VerifyResponse := VerifyResponse{}
	json.Unmarshal(jsonString, &VerifyResponse)
	return &VerifyResponse
}
