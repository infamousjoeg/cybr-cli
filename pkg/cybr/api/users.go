package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/queries"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/responses"
	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// UnsuspendUser activates a suspended user. It does not activate an inactive user.
func (c Client) UnsuspendUser(username string) error {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Users/%s", c.BaseURL, url.QueryEscape(username))

	body := requests.UnsuspendUser{
		Suspended: false,
	}

	response, err := httpJson.Put(url, c.SessionToken, body, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to unsuspend user '%s'. %s. %s", username, string(returnedError), err)
	}
	return nil
}

// ListUsers returns a list of all existing users in the Vault except for the Master and the Batch built-in users.
func (c Client) ListUsers(query *queries.ListUsers) (responses.ListUsers, error) {
	url := fmt.Sprintf("%s/PasswordVault/api/Users%s", c.BaseURL, httpJson.GetURLQuery(query))

	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return responses.ListUsers{}, fmt.Errorf("Failed to list users. %s. %s", string(returnedError), err)
	}

	jsonString, _ := json.Marshal(response)
	ListUsersResponse := responses.ListUsers{}
	err = json.Unmarshal(jsonString, &ListUsersResponse)
	return ListUsersResponse, err
}

// DeleteUser from PAS
func (c Client) DeleteUser(userID int) error {
	url := fmt.Sprintf("%s/PasswordVault/api/Users/%d", c.BaseURL, userID)

	response, err := httpJson.Delete(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to delete user '%d'. %s. %s", userID, string(returnedError), err)
	}

	return nil
}
