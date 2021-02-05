package api_test

import (
	"strings"
	"testing"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
)

func TestListSafesSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	safes, err := client.ListSafes()
	if err != nil {
		t.Errorf("Failed to list safes. %s", err)
	}

	if len(safes.Safes) == 0 {
		t.Errorf("At least one safe should be listed. %s", err)
	}
}

func TestListSafeMembersSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	members, err := client.ListSafeMembers("PasswordManager")
	if err != nil {
		t.Errorf("Failed to list safes. %s", err)
	}

	if len(members.Members) == 0 {
		t.Errorf("At least one member should be listed. %s", err)
	}
}

func TestListSafeMembersInvalidSafeName(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	_, err = client.ListSafeMembers("notReal")
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("Expecting 404 error to be returned but got. %s", err)
	}
}

func TestAddRemoveSafeSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	newSafe := requests.AddSafe{
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
	client, err := defaultPASAPIClient(t)

	safeName := "notRealSafeName"
	err = client.DeleteSafe(safeName)
	if err == nil {
		t.Errorf("Client returned successful safe deletion even though safe '%s' should not exist", safeName)
	}
}

func TestAddRemoveSafeMemberSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	safeName := "PasswordManager"
	memberName := "test-add-member"

	retrieveAccounts := requests.PermissionKeyValue{
		Key:   "RetrieveAccounts",
		Value: true,
	}

	addMember := requests.AddSafeMember{
		Member: requests.AddSafeMemberInternal{
			MemberName:  memberName,
			SearchIn:    "Vault",
			Permissions: []requests.PermissionKeyValue{retrieveAccounts},
		},
	}

	err = client.AddSafeMember(safeName, addMember)
	if err != nil {
		t.Errorf("Failed to add member to safe. %s", err)
	}

	err = client.RemoveSafeMember(safeName, memberName)
	if err != nil {
		t.Errorf("Failed to remove member from safe. %s", err)
	}
}

func TestAddMemberInvalidMemberName(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	safeName := "PasswordManager"
	memberName := "notReal"

	retrieveAccounts := requests.PermissionKeyValue{
		Key:   "RetrieveAccounts",
		Value: true,
	}

	addMember := requests.AddSafeMember{
		Member: requests.AddSafeMemberInternal{
			MemberName:  memberName,
			SearchIn:    "Vault",
			Permissions: []requests.PermissionKeyValue{retrieveAccounts},
		},
	}

	err = client.AddSafeMember(safeName, addMember)
	if err == nil {
		t.Errorf("Added a non-existent member. This should not happen")
	}
}

func TestRemoveMemberInvalidMemberName(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	safeName := "PasswordManager"
	memberName := "notReal"

	err = client.RemoveSafeMember(safeName, memberName)
	if err == nil {
		t.Errorf("Removed a non-existent member. This should not happen")
	}
}
