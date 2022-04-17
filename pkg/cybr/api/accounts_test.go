package api_test

import (
	"testing"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/queries"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
)

var (
	accountSafeName  = "CLI_ACCOUNTS_TEST"
	accountID        = "110_3"
	accountSSHKeyID  = "110_15"
	invalidAccountID = "202_5"
)

func TestListAccountSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	query := &queries.ListAccounts{}

	accounts, err := client.ListAccounts(query)
	if err != nil {
		t.Errorf("Failed to list accounts. %s", err)
	}

	if accounts.Count == 0 {
		t.Error("At least one account should be returned")
	}
}

func TestListAccountSearchSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	query := &queries.ListAccounts{
		Search: accountSafeName,
	}

	_, err = client.ListAccounts(query)
	if err != nil {
		t.Errorf("Failed to list accounts. %s", err)
	}
}

func TestGetAccountSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	account, err := client.GetAccount("110_3")
	if err != nil {
		t.Errorf("Failed to get account. %s", err)
	}

	if account.UserName != "test" {
		t.Errorf("Retrieved invalid account. Account has username '%s' and should be 'test'", account.UserName)
	}
}

func TestGetAccountInvalidAccountID(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	_, err = client.GetAccount("not_good")
	if err == nil {
		t.Errorf("The account does not exists but getting account was successful. This should not happen")
	}
}

func TestAddAndDeleteAccountSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	account := requests.AddAccount{
		SafeName:   accountSafeName,
		Address:    "10.0.0.1",
		UserName:   "add-account-test",
		PlatformID: "UnixSSH",
		SecretType: "password",
		// file deepcode ignore HardcodedPassword/test: dummy secret for testing
		Secret: "superSecret",
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
	client, err := defaultPASAPIClient(t)

	account := requests.AddAccount{
		SafeName:   "invalidSafeName",
		Address:    "10.0.0.1",
		UserName:   "add-account-test",
		PlatformID: "UnixSSH",
		SecretType: "password",
		// file deepcode ignore HardcodedPassword/test: dummy secret for testing
		Secret: "superSecret",
	}

	_, err = client.AddAccount(account)
	if err == nil {
		t.Errorf("Added account to invalid safe name. This should not happen")
	}
}

func TestDeleteAccountInvalidAccountID(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	err = client.DeleteAccount("invalid_ID")
	if err == nil {
		t.Errorf("Delete account but it does not exists. This should not happen")
	}
}

func TestGetAccountPasswordSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	result, err := client.GetAccountPassword(accountID, requests.GetAccountPassword{})
	if err != nil {
		t.Errorf("Failed to get account password. %s", err)
	}

	if result != "Cyberark1" {
		t.Errorf(result)
		t.Errorf("Password retrieved is invalid")
	}
}

func TestGetAccountPasswordInvalidAccountID(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	_, err = client.GetAccountPassword(invalidAccountID, requests.GetAccountPassword{})
	if err == nil {
		t.Errorf("Got password but should have failed")
	}
}

func TestGetAccountSSHKeySuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	_, err = client.GetAccountSSHKey(accountSSHKeyID, requests.GetAccountPassword{})
	if err != nil {
		t.Errorf("Failed to get account password. %s", err)
	}
}

func TestGetAccountSSHKeyInvalidAccountID(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	_, err = client.GetAccountSSHKey(invalidAccountID, requests.GetAccountPassword{})
	if err == nil {
		t.Errorf("Got password but should have failed")
	}
}

func TestVerifyAccountCredentialsSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	err = client.VerifyAccountCredentials(accountID)
	if err != nil {
		t.Errorf("Failed to get account password. %s", err)
	}
}

func TestVerifyAccountCredentialsInvalidAccount(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	err = client.VerifyAccountCredentials(invalidAccountID)
	if err == nil {
		t.Errorf("Set account for verify but it should not exist")
	}
}

func TestChangeAccountCredentialsSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	err = client.ChangeAccountCredentials(accountID, false)
	if err != nil {
		t.Errorf("Failed to get account password. %s", err)
	}
}

func TestChangeAccountCredentialsInvalidAccount(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	err = client.ChangeAccountCredentials(invalidAccountID, false)
	if err == nil {
		t.Errorf("Set account for change but it should not exist")
	}
}

func TestReconileAccountCredentialsSuccess(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	err = client.ReconileAccountCredentials(accountID)
	if err != nil {
		t.Errorf("Failed to get account password. %s", err)
	}
}

func TestReconcileAccountCredentialsInvalidAccount(t *testing.T) {
	client, err := defaultPASAPIClient(t)

	err = client.ReconileAccountCredentials(invalidAccountID)
	if err == nil {
		t.Errorf("Set account for reconcile but it should not exist")
	}
}
