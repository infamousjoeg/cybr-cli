package api_test

import (
	"testing"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
)

func TestUnsuspendUserSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	err = client.UnsuspendUser(username)
	if err != nil {
		t.Errorf("Failed to unsuspend user '%s'. %s", username, err)
	}
}

func TestUnsuspendUserInvalidUsername(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	invalidUsername := "invalidUsername"
	err = client.UnsuspendUser(invalidUsername)
	if err == nil {
		t.Errorf("Unsuspended user '%s' but user should not exist. This should never happen", invalidUsername)
	}
}

func TestListUsersSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)
	query := &api.ListUsersQueryParams{}

	_, err = client.ListUsers(query)
	if err != nil {
		t.Errorf("Failed to list users. %s", err)
	}
}

func TestListUsersInvalidUsername(t *testing.T) {
	client, err := defaultPASAPIClient(t)
	query := &api.ListUsersQueryParams{
		Search: "invalidUsername",
	}

	users, err := client.ListUsers(query)
	if err != nil {
		t.Errorf("Failed to list users. %s", err)
	}

	if users.Total > 0 {
		t.Errorf("No users should be returned")
	}
}

func TestAddRemoveUserSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)
	newUser := requests.AddUser{
		Username:              "testAddRemove",
		InitialPassword:       "Cyberark1",
		UserType:              "EPVUser",
		AuthenticationMethod:  []string{"AuthTypePass"},
		Location:              "\\\\",
		EnableUser:            true,
		ChangePassOnNextLogon: false,
		PasswordNeverExpires:  true,
		Description:           "testAddRemove user",
		DistinguishedName:     "testAddRemove@cyberark",
		PersonalDetails: requests.AddUserPersonalDetails{
			Department: "R&D",
		},
	}

	response, err := client.AddUser(newUser)
	if err != nil {
		t.Errorf("%s", err)
	}

	err = client.DeleteUser(response.ID)
	if err != nil {
		t.Errorf("Failed to delete user '%d'", response.ID)
	}
}
