package api_test

import (
	"testing"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
)

var (
	accountSafeName = "CLI_ACCOUNTS_TEST"
)

func TestListAccountSuccess(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "ldap",
	}

	creds := pasapi.LogonRequest{
		Username: username,
		Password: password,
	}

	err := client.Logon(creds)
	if err != nil {
		t.Errorf("Failed to logon. %s", err)
	}

	query := &pasapi.ListAccountQueryParams{}

	accounts, err := client.ListAccounts(query)
	if err != nil {
		t.Errorf("Failed to list accounts. %s", err)
	}

	if accounts.Count == 0 {
		t.Error("At least one account should be returned")
	}
}

func TestListAccountSearchSuccess(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "ldap",
	}

	creds := pasapi.LogonRequest{
		Username: username,
		Password: password,
	}

	err := client.Logon(creds)
	if err != nil {
		t.Errorf("Failed to logon. %s", err)
	}

	query := &pasapi.ListAccountQueryParams{
		Search: accountSafeName,
	}

	accounts, err := client.ListAccounts(query)
	if err != nil {
		t.Errorf("Failed to list accounts. %s", err)
	}

	if accounts.Count != 1 {
		t.Error("Only 1 account should be returned")
	}
}

func TestGetAccountSuccess(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "ldap",
	}

	creds := pasapi.LogonRequest{
		Username: username,
		Password: password,
	}

	err := client.Logon(creds)
	if err != nil {
		t.Errorf("Failed to logon. %s", err)
	}

	account, err := client.GetAccount("110_3")
	if err != nil {
		t.Errorf("Failed to get account. %s", err)
	}

	if account.UserName != "test" {
		t.Errorf("Retrieved invalid account. Account has username '%s' and should be 'test'", account.UserName)
	}
}

func TestGetAccountInvalidAccountID(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "ldap",
	}

	creds := pasapi.LogonRequest{
		Username: username,
		Password: password,
	}

	err := client.Logon(creds)
	if err != nil {
		t.Errorf("Failed to logon. %s", err)
	}

	_, err = client.GetAccount("not_good")
	if err == nil {
		t.Errorf("The account does not exists but getting account was successful. This should not happen")
	}
}

func TestAddAndDeleteAccountSuccess(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "ldap",
	}

	creds := pasapi.LogonRequest{
		Username: username,
		Password: password,
	}

	err := client.Logon(creds)
	if err != nil {
		t.Errorf("Failed to logon. %s", err)
	}

	account := pasapi.AddAccountRequest{
		SafeName:   accountSafeName,
		Address:    "10.0.0.1",
		UserName:   "add-account-test",
		PlatformID: "UnixSSH",
		SecretType: "password",
		Secret:     "superSecret",
	}

	addedAccount, err := client.AddAccount(account)
	if err != nil {
		t.Errorf("Failed to add account. %s", err)
	}

	if addedAccount.UserName != account.UserName {
		t.Errorf("The added account has a different username than the provided account. This should not occur")
	}

	err = client.DeleteAccount(addedAccount.ID)
	if err != nil {
		t.Errorf("Failed to delete account even though it should exists. %s", err)
	}
}

func TestAddAccountInvalidSafeName(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "ldap",
	}

	creds := pasapi.LogonRequest{
		Username: username,
		Password: password,
	}

	err := client.Logon(creds)
	if err != nil {
		t.Errorf("Failed to logon. %s", err)
	}

	account := pasapi.AddAccountRequest{
		SafeName:   "invalidSafeName",
		Address:    "10.0.0.1",
		UserName:   "add-account-test",
		PlatformID: "UnixSSH",
		SecretType: "password",
		Secret:     "superSecret",
	}

	_, err = client.AddAccount(account)
	if err == nil {
		t.Errorf("Added account to invalid safe name. This should not happen")
	}
}

func TestDeleteAccountInvalidAccountID(t *testing.T) {
	client := pasapi.Client{
		BaseURL:  hostname,
		AuthType: "ldap",
	}

	creds := pasapi.LogonRequest{
		Username: username,
		Password: password,
	}

	err := client.Logon(creds)
	if err != nil {
		t.Errorf("Failed to logon. %s", err)
	}

	err = client.DeleteAccount("invalid_ID")
	if err == nil {
		t.Errorf("Delete account but it does not exists. This should not happen")
	}
}
