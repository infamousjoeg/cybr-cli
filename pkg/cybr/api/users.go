package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/responses"
	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// UnsuspendUserRequest request used when unsuspending user
type UnsuspendUserRequest struct {
	Suspended bool `json:"Suspended"`
}

// ListUsersQueryParams represents valid query parameters when listing users
type ListUsersQueryParams struct {
	Search string `query_key:"search"`
	Filter string `query_key:"filter"`
}

// ListUsersResponse response when listing users
type ListUsersResponse struct {
	Users []UserResponse `json:"Users"`
	Total int            `json:"Total"`
}

// UserResponse represents one user in ListUsersResponse
type UserResponse struct {
	ID                 int                     `json:"id"`
	Username           string                  `json:"username"`
	Source             string                  `json:"source"`
	UserType           string                  `json:"userType"`
	ComponentUser      bool                    `json:"componentUser"`
	VaultAuthorization []string                `json:"vaultAuthorization"`
	Location           string                  `json:"location"`
	PersonalDetails    PersonalDetailsResponse `json:"personalDetails"`
}

// PersonalDetailsResponse represents one users personal details
type PersonalDetailsResponse struct {
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
}

// UnsuspendUser activates a suspended user. It does not activate an inactive user.
func (c Client) UnsuspendUser(username string) error {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Users/%s", c.BaseURL, url.QueryEscape(username))

	body := UnsuspendUserRequest{
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
func (c Client) ListUsers(query *ListUsersQueryParams) (ListUsersResponse, error) {
	url := fmt.Sprintf("%s/PasswordVault/api/Users%s", c.BaseURL, httpJson.GetURLQuery(query))

	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return ListUsersResponse{}, fmt.Errorf("Failed to list users. %s. %s", string(returnedError), err)
	}

	jsonString, _ := json.Marshal(response)
	ListUsersResponse := ListUsersResponse{}
	err = json.Unmarshal(jsonString, &ListUsersResponse)
	return ListUsersResponse, err
}

// AddUser to PAS
func (c Client) AddUser(user requests.AddUser) (responses.AddUser, error) {
	url := fmt.Sprintf("%s/PasswordVault/api/Users", c.BaseURL)

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
