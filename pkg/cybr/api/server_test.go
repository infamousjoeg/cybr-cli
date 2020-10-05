package api_test

import (
	"testing"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
)

func TestServerVerifySuccess(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "cyberark",
	}

	response, err := client.ServerVerify()
	if err != nil {
		t.Errorf("Failed to get Server verify info from PVWA. %s", err)
	}

	if response.ServerName != "Vault" {
		t.Errorf("Invalid server name '%s'", response.ServerName)
	}
}

func TestServerVerifyInvalidHostName(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  "https://invalidHostName",
		AuthType: "cyberark",
	}

	_, err := client.ServerVerify()
	if err == nil {
		t.Errorf("Successfully got Server Verify info but shouldn't have. %s", err)
	}
}
