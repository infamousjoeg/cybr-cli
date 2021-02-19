package api_test

import (
	"testing"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/queries"
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
	query := &queries.ListUsers{}

	_, err = client.ListUsers(query)
	if err != nil {
		t.Errorf("Failed to list users. %s", err)
	}
}

func TestListUsersInvalidUsername(t *testing.T) {
	client, err := defaultPASAPIClient(t)
	query := &queries.ListUsers{
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

func TestRemoveUserInvalidUserID(t *testing.T) {
	client, err := defaultPASAPIClient(t)
	err = client.DeleteUser(100000)
	if err == nil {
		t.Errorf("Error should have been returned becauste UserID is invalid")
	}
}

func TestAddRemoveUser(t *testing.T) {
	client, err := defaultPASAPIClient(t)
	user := requests.AddUser{
		Username:               "testcliuser",
		UserType:               "EPVUser",
		InitialPassword:        "Cyberark1",
		AuthenticationMethod:   []string{"AuthTypePass"},
		Location:               "\\",
		UnAuthorizedInterfaces: []string{"PSM", "PSMP"},
		VaultAuthorization:     []string{"AddSafes", "AuditUsers"},
		EnableUser:             true,
		ChangePassOnNextLogon:  false,
		PasswordNeverExpires:   true,
		Description:            "test cli user",
	}

	addedUser, err := client.AddUser(user)
	if err != nil {
		t.Errorf("Failed to add user '%s'. %s", user.Username, err)
	}

	err = client.DeleteUser(addedUser.ID)
	if err != nil {
		t.Errorf("Failed to delete user '%d'. %s", addedUser.ID, err)
	}
}
