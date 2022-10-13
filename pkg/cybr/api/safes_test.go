package api_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/queries"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
)

// The content will look like
// port=something, sp
func keyValueStringToMap(content string) (map[string]string, error) {
	if content == "" {
		return nil, nil
	}

	if !strings.Contains(content, "=") {
		return nil, fmt.Errorf("Invalid platform prop content. The provided content does not container a '='")
	}

	m := make(map[string]string)

	// TODO: Gotta be a better way to do this
	replaceWith := "^||||^"

	// If the address or property contains a `\,` then replace
	content = strings.ReplaceAll(content, "\\,", replaceWith)
	props := strings.Split(content, ",")
	for _, prop := range props {
		if !strings.Contains(prop, "=") {
			return nil, fmt.Errorf("Property '%s' is invalid because it does not contain a '=' to seperate key from value", prop)
		}
		kvs := strings.SplitN(prop, "=", 2)
		key := strings.Trim(kvs[0], " ")
		value := strings.Trim(strings.ReplaceAll(kvs[1], replaceWith, ","), " ")
		m[key] = value
	}

	return m, nil
}

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

	query := &queries.ListSafeMembers{
		Search: "",
		Sort:   "",
		Offset: 0,
		Limit:  50,
		Filter: "",
	}

	members, err := client.ListSafeMembers("PasswordManager", query)
	if err != nil {
		t.Errorf("Failed to list safes. %s", err)
	}

	if len(members.Members) == 0 {
		t.Errorf("At least one member should be listed. %s", err)
	}
}

func TestListSafeMembersInvalidSafeName(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	query := &queries.ListSafeMembers{
		Search: "",
		Sort:   "",
		Offset: 0,
		Limit:  50,
		Filter: "",
	}

	_, err = client.ListSafeMembers("notReal", query)
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

	retrieveAccounts, err := keyValueStringToMap("RetrieveAccounts=true")
	if err != nil {
		t.Errorf("Failed to parse props to map. %s", err)
	}

	addMember := requests.AddSafeMember{
		MemberName:  memberName,
		SearchIn:    "Vault",
		Permissions: retrieveAccounts,
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

	retrieveAccounts, err := keyValueStringToMap("RetrieveAccounts=true")
	if err != nil {
		t.Errorf("Failed to parse props to map. %s", err)
	}

	addMember := requests.AddSafeMember{
		MemberName:  memberName,
		SearchIn:    "Vault",
		Permissions: retrieveAccounts,
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
