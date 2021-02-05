package api_test

import (
	"strings"
	"testing"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
)

func TestCyberarkLogonSuccess(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "cyberark",
	}

	creds := requests.Logon{
		Username: username,
		Password: password,
	}

	err := client.Logon(creds)
	if err != nil {
		t.Errorf("Failed to logon. %s", err)
	}
}

func TestCyberarkLogonInvalidCreds(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "cyberark",
	}

	creds := requests.Logon{
		Username: "notReal",
		Password: password,
	}

	err := client.Logon(creds)
	if err == nil {
		t.Errorf("Successfully logged in but shouldn't have. %s", err)
	}
}

func TestCyberarkLogonInvalidHostName(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  "https://invalidhostname",
		AuthType: "cyberark",
	}

	creds := requests.Logon{
		Username: "notReal",
		Password: password,
	}

	err := client.Logon(creds)
	if err == nil {
		t.Errorf("Successfully logged in but shouldn't have. %s", err)
	}
}

func TestLogonInvalidAuthType(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "notGood",
	}

	creds := requests.Logon{
		Username: username,
		Password: password,
	}

	err := client.Logon(creds)
	if err == nil {
		t.Errorf("Successfully logged in but shouldn't have. %s", err)
	}

	if !strings.Contains(err.Error(), "Invalid auth type 'notGood'") {
		t.Errorf("Recieved incorrect error message. %s", err)
	}
}

func TestCyberarkLogoffSuccess(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "cyberark",
	}

	creds := requests.Logon{
		Username: username,
		Password: password,
	}

	err := client.Logon(creds)
	if err != nil {
		t.Errorf("Failed to logon. %s", err)
	}

	err = client.Logoff()
	if err != nil {
		t.Errorf("Failed to logoff. %s", err)
	}
}

func TestCyberarkLogoffFailNotLoggedIn(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "cyberark",
	}

	err := client.Logoff()
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("Expected to recieve 401 statuc code. %s", err)
	}
}
