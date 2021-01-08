package api_test

import (
	"strings"
	"testing"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
)

func TestListSafesSuccess(t *testing.T) {
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

	safes, err := client.ListSafes()
	if err != nil {
		t.Errorf("Failed to list safes. %s", err)
	}

	if len(safes.Safes) == 0 {
		t.Errorf("At least one safe should be listed. %s", err)
	}
}

func TestListSafeMembersSuccess(t *testing.T) {
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

	members, err := client.ListSafeMembers("PasswordManager")
	if err != nil {
		t.Errorf("Failed to list safes. %s", err)
	}

	if len(members.Members) == 0 {
		t.Errorf("At least one member should be listed. %s", err)
	}
}

func TestListSafeMembersInvalidSafeName(t *testing.T) {
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

	_, err = client.ListSafeMembers("notReal")
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("Expecting 404 error to be returned but got. %s", err)
	}
}

func TestAddRemoveSafeSuccess(t *testing.T) {
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
	newSafe := pasapi.AddSafeRequest{
		SafeName:              "TestCreateDelete",
		Description:           "Testing creating and deleteing a safe",
		OLACEnabled:           false,
		ManagingCPM:           "PasswordManager",
		NumberOfDaysRetention: 0,
	}
	err = client.AddSafe(newSafe)
	if err != nil {
		t.Errorf("Failed to create safe '%s' even though it should have been created successfully. %s", newSafe.SafeName, err)
	}

	err = client.DeleteSafe(newSafe.SafeName)
	if err != nil {
		t.Errorf("Failed to delete safe '%s' even though it should exist and should be deletable. %s", newSafe.SafeName, err)
	}
}

func TestRemoveSafeFail(t *testing.T) {
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

	safeName := "notRealSafeName"
	err = client.DeleteSafe(safeName)
	if err == nil {
		t.Errorf("Client returned successful safe deletion even though safe '%s' should not exist", safeName)
	}
}
