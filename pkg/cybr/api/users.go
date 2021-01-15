package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// UnsuspendUserRequest request used when unsuspending user
type UnsuspendUserRequest struct {
	Suspended bool `json:"Suspended"`
}

// UnsuspendUser activates a suspended user. It does not activate an inactive user.
func (c Client) UnsuspendUser(username string) error {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Users/%s", c.BaseURL, url.QueryEscape(username))

	body := UnsuspendUserRequest{
		Suspended: false,
	}

	response, err := httpJson.Put(url, c.SessionToken, body, c.InsecureTLS)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to unsuspend user '%s'. %s. %s", username, string(returnedError), err)
	}
	return nil
}
