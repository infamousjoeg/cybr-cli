package api_test

import (
	"strings"
	"testing"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
)

func TestListApplicationSuccess(t *testing.T) {
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

	apps, err := client.ListApplications("\\")
	if err != nil {
		t.Errorf("Failed to list applications. %s", err)
	}

	if len(apps.Application) == 0 {
		t.Error("At least one application should be returned")
	}
}

func TestListApplicationInvalidLocation(t *testing.T) {
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

	_, err = client.ListApplications("notReal")
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("Expecting 400 status code. %s", err)
	}
}

func TestListApplicationAuthenticationMethodsSuccess(t *testing.T) {
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

	_, err = client.ListApplicationAuthenticationMethods("test")
	if err != nil {
		t.Errorf("Failed to list application authentication methods. %s", err)
	}
}

func TestListApplicationAuthenticationMethodsInvalidApplicationID(t *testing.T) {
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

	_, err = client.ListApplicationAuthenticationMethods("notReal")
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("Expecting 404 status code. %s", err)
	}
}
