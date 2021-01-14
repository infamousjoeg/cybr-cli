package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// UnsuspendUser activates a suspended user. It does not activate an inactive user.
func (c Client) UnsuspendUser(username string) error {
	url := fmt.Sprintf("%s/PasswordVault/api/Users/%s/Activate", c.BaseURL, url.QueryEscape(username))
	response, err := httpJson.Post(url, c.SessionToken, nil, c.InsecureTLS)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to unsuspend user '%s'. %s. %s", username, string(returnedError), err)
	}
	return nil
}
