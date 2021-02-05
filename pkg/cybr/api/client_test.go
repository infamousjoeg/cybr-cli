package api_test

import (
	"os"
	"testing"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
)

var (
	hostname = os.Getenv("PAS_HOSTNAME")
	username = os.Getenv("PAS_USERNAME")
	password = os.Getenv("PAS_PASSWORD")
)

func defaultPASAPIClient(t *testing.T) (pasapi.Client, error) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "ldap",
	}

	creds := requests.Logon{
		Username: username,
		Password: password,
	}

	err := client.Logon(creds)
	if err != nil {
		t.Errorf("Failed to logon. %s", err)
	}
	return client, err
}

func TestIsValidSuccess(t *testing.T) {
	validAuthTypes := []string{"cyberark", "ldap"}

	for _, validAuthType := range validAuthTypes {
		client := pasapi.Client{
			AuthType: validAuthType,
		}

		err := client.IsValid()
		if err != nil {
			t.Errorf("Auth type '%s' is valid however client is returning not valid. %s", validAuthType, err)
		}
	}
}

func TestIsValidFail(t *testing.T) {
	validAuthTypes := []string{"invalidAuthType", "", "123456"}

	for _, validAuthType := range validAuthTypes {
		client := pasapi.Client{
			AuthType: validAuthType,
		}

		err := client.IsValid()
		if err == nil {
			t.Errorf("Auth type '%s' is not valid however client is returning valid. %s", validAuthType, err)
		}
	}
}

func TestSetGetRemoveConfigSuccess(t *testing.T) {
	// I tested these three functions together so we can create, read and delete
	// without worrying about the order of test execution
	client := pasapi.Client{
		BaseURL:      hostname,
		AuthType:     "cyberark",
		InsecureTLS:  false,
		SessionToken: "",
	}

	err := client.SetConfig()
	if err != nil {
		t.Errorf("Failed to set the pasapi configuration. %s", err)
	}

	client, err = pasapi.GetConfig()
	if err != nil {
		t.Errorf("Failed to retrieve pasapi config. %s", err)
	}

	err = client.RemoveConfig()
	if err != nil {
		t.Errorf("Failed to remove pasapi config. %s", err)
	}
}

func TestGetConfigFail(t *testing.T) {
	_, err := pasapi.GetConfig()
	if err == nil {
		t.Errorf("Successfully retrieved pasapi config but pasapi config should not exist")
	}
}
