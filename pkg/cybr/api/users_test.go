package api_test

import (
	"testing"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/queries"
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
