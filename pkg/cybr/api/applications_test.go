package api_test

import (
	"strings"
	"testing"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
)

func TestListApplicationSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	apps, err := client.ListApplications("\\")
	if err != nil {
		t.Errorf("Failed to list applications. %s", err)
	}

	if len(apps.Application) == 0 {
		t.Error("At least one application should be returned")
	}
}

func TestListApplicationInvalidLocation(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	_, err = client.ListApplications("notReal")
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("Expecting 400 status code. %s", err)
	}
}

func TestListApplicationAuthenticationMethodsSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	_, err = client.ListApplicationAuthenticationMethods("test-list-authn")
	if err != nil {
		t.Errorf("Failed to list application authentication methods. %s", err)
	}
}

func TestListApplicationAuthenticationMethodsInvalidApplicationID(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	_, err = client.ListApplicationAuthenticationMethods("notReal")
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("Expecting 404 status code. %s", err)
	}
}

func TestAddDeleteApplicationSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	newApplication := requests.AddApplication{
		Application: requests.Application{
			AppID:               "test-api-app1",
			Description:         "Some type of description",
			Location:            "\\",
			AccessPermittedFrom: 0,
			AccessPermittedTo:   23,
		},
	}

	err = client.AddApplication(newApplication)
	if err != nil {
		t.Errorf("Failed to add application. %s", err)
	}

	err = client.DeleteApplication(newApplication.Application.AppID)
	if err != nil {
		t.Errorf("Failed to delete application. %s", err)
	}
}

func TestDeleteApplicationInvalidAppID(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	err = client.DeleteApplication("invalid-app-id")
	if err == nil {
		t.Errorf("Delete application but it should not exist. This should not happen")
	}
}

func TestAddDeleteApplicationAuthenticationMethodsSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	newApplication := requests.AddApplication{
		Application: requests.Application{
			AppID:               "test-api-method-app1",
			Description:         "Some type of description",
			Location:            "\\",
			AccessPermittedFrom: 0,
			AccessPermittedTo:   23,
		},
	}

	err = client.AddApplication(newApplication)
	if err != nil {
		t.Errorf("Failed to add application. %s", err)
	}

	newAuthnMethod := requests.AddApplicationAuthentication{
		Authentication: requests.ApplicationAuthenticationMethod{
			AuthType:             "path",
			AuthValue:            "/some/path",
			IsFolder:             false,
			AllowInternalScripts: false,
		},
	}

	err = client.AddApplicationAuthenticationMethod(newApplication.Application.AppID, newAuthnMethod)
	if err != nil {
		t.Errorf("Failed to add application authentication method. %s", err)
	}

	err = client.DeleteApplicationAuthenticationMethod(newApplication.Application.AppID, "1")
	if err != nil {
		t.Errorf("Failed to delete application authentication method. %s", err)
	}

	err = client.DeleteApplication(newApplication.Application.AppID)
	if err != nil {
		t.Errorf("Failed to delete application. %s", err)
	}
}
