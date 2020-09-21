package api

import (
	"encoding/json"
	"fmt"

	httpJson "github.com/infamousjoeg/pas-api-go/pkg/helpers"
)

// ServerVerify is an unauthenticated endpoint for testing Web Service availability
func ServerVerify(hostname string) (string, error) {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Verify", hostname)
	response, err := httpJson.SendRequest(url, "GET", "", "")
	// Marshal (convert) returned map string interface to JSON
	resJSON, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("Unable to marshal map to JSON. %s", err)
	}
	return string(resJSON), err
}
