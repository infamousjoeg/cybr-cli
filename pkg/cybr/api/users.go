package api

import (
	"encoding/json"
	"fmt"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/queries"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/responses"
	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// UnsuspendUser activates a suspended user. It does not activate an inactive user.
func (c Client) UnsuspendUser(userID int) error {
	url := fmt.Sprintf("%s/passwordvault/api/Users/%d/activate", c.BaseURL, userID)

	response, err := httpJson.Post(url, c.SessionToken, nil, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to unsuspend user with id '%d'. %s. %s", userID, string(returnedError), err)
	}
	return nil
}

// ListUsers returns a list of all existing users in the Vault except for the Master and the Batch built-in users.
func (c Client) ListUsers(query *queries.ListUsers) (responses.ListUsers, error) {
	url := fmt.Sprintf("%s/passwordvault/api/Users%s", c.BaseURL, httpJson.GetURLQuery(query))

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
	url := fmt.Sprintf("%s/passwordvault/api/Users/%d", c.BaseURL, userID)

	response, err := httpJson.Delete(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to delete user '%d'. %s. %s", userID, string(returnedError), err)
	}

	return nil
}

// AddUser to PAS
func (c Client) AddUser(user requests.AddUser) (responses.AddUser, error) {
	url := fmt.Sprintf("%s/passwordvault/api/Users", c.BaseURL)

	response, err := httpJson.Post(url, c.SessionToken, user, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return responses.AddUser{}, fmt.Errorf("Failed to add user '%s'. %s. %s", user.Username, string(returnedError), err)
	}

	jsonString, _ := json.Marshal(response)
	addUserResponse := responses.AddUser{}
	err = json.Unmarshal(jsonString, &addUserResponse)
	return addUserResponse, err
}
