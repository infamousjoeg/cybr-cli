package api_test

import (
	"testing"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
)

func TestUnsuspendUserSuccess(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "cyberark",
	}

	creds := pasapi.LogonRequest{
		Username: username,
		Password: password,
	}

	err := client.Logon(creds)
	if err != nil {
		t.Errorf("Failed to logon. %s", err)
	}

	err = client.UnsuspendUser(username)
	if err != nil {
		t.Errorf("Failed to unsuspend user '%s'. %s", username, err)
	}
}

func TestUnsuspendUserInvalidUsername(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "cyberark",
	}

	creds := pasapi.LogonRequest{
		Username: username,
		Password: password,
	}

	err := client.Logon(creds)
	if err != nil {
		t.Errorf("Failed to logon. %s", err)
	}

	invalidUsername := "invalidUsername"
	err = client.UnsuspendUser(invalidUsername)
	if err == nil {
		t.Errorf("Unsuspended user '%s' but user should not exist. This should never happen", invalidUsername)
	}
}
