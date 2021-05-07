package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/responses"
	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// ServerVerify is an unauthenticated endpoint for testing Web Service availability
func (c Client) ServerVerify() (*responses.ServerVerify, error) {
	ctx := context.TODO()
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Verify", c.BaseURL)
	response, err := httpJson.SendRequest(ctx, url, "GET", "", nil, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.ServerVerify{}, fmt.Errorf("Error verifying PAS REST API Web Service. %s", err)
	}
	jsonString, _ := json.Marshal(response)
	VerifyResponse := responses.ServerVerify{}
	err = json.Unmarshal(jsonString, &VerifyResponse)
	return &VerifyResponse, err
}
