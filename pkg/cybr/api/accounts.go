package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/queries"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/responses"
	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// ListAccounts CyberArk user has access to
func (c Client) ListAccounts(query *queries.ListAccounts) (*responses.ListAccount, error) {
	url := fmt.Sprintf("%s/passwordvault/api/Accounts%s", c.BaseURL, httpJson.GetURLQuery(query))
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.ListAccount{}, fmt.Errorf("Failed to list accounts. %s", err)
	}

	jsonString, _ := json.Marshal(response)
	ListAccountsResponse := responses.ListAccount{}
	err = json.Unmarshal(jsonString, &ListAccountsResponse)
	return &ListAccountsResponse, err
}

// GetAccount details for specific account
func (c Client) GetAccount(accountID string) (*responses.GetAccount, error) {
	url := fmt.Sprintf("%s/passwordvault/api/Accounts/%s", c.BaseURL, accountID)
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.GetAccount{}, fmt.Errorf("Failed to get account. %s", err)
	}

	jsonString, _ := json.Marshal(response)
	GetAccountResponse := &responses.GetAccount{}
	err = json.Unmarshal(jsonString, GetAccountResponse)
	return GetAccountResponse, err
}

// AddAccount to cyberark
func (c Client) AddAccount(account requests.AddAccount) (*responses.GetAccount, error) {
	url := fmt.Sprintf("%s/passwordvault/api/Accounts", c.BaseURL)
	logger := c.GetLogger().AddSecret(account.Secret)
	response, err := httpJson.Post(url, c.SessionToken, account, c.InsecureTLS, logger)
	logger = logger.ClearSecrets()

	if err != nil {
		returnedError, _ := json.Marshal(response)
		return &responses.GetAccount{}, fmt.Errorf("Failed to add account. %s. %s", string(returnedError), err)
	}

	jsonString, _ := json.Marshal(response)
	GetAccountResponse := &responses.GetAccount{}
	err = json.Unmarshal(jsonString, GetAccountResponse)
	return GetAccountResponse, err
}

// DeleteAccount from cyberark
func (c Client) DeleteAccount(accountID string) error {
	url := fmt.Sprintf("%s/passwordvault/api/Accounts/%s", c.BaseURL, accountID)
	response, err := httpJson.Delete(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to delete account '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return nil
}

// GetJITAccess from a specific account
func (c Client) GetJITAccess(accountID string) error {
	url := fmt.Sprintf("%s/passwordvault/api/Accounts/%s/grantAdministrativeAccess", c.BaseURL, accountID)
	response, err := httpJson.Post(url, c.SessionToken, nil, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to get JIT access for account '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return nil
}

// RevokeJITAccess from a specific account
func (c Client) RevokeJITAccess(accountID string) error {
	url := fmt.Sprintf("%s/passwordvault/api/Accounts/%s/RevokeAdministrativeAccess", c.BaseURL, accountID)
	response, err := httpJson.Post(url, c.SessionToken, nil, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to revoke JIT access for account '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return nil
}

// GetAccountPassword This method enables users to retrieve the password or SSH key of an existing account that is identified by its Account ID. It enables users to specify a reason and ticket ID, if required
func (c Client) GetAccountPassword(accountID string, request requests.GetAccountPassword) (string, error) {
	url := fmt.Sprintf("%s/passwordvault/api/Accounts/%s/Password/Retrieve", c.BaseURL, accountID)

	response, err := httpJson.SendRequestRaw(url, "POST", c.SessionToken, request, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return "", fmt.Errorf("Failed to retrieve the account password '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return strings.Trim(string(response), "\""), nil
}

// GetAccountSSHKey This method enables users to retrieve the password or SSH key of an existing account that is identified by its Account ID. It enables users to specify a reason and ticket ID, if required
func (c Client) GetAccountSSHKey(accountID string, request requests.GetAccountPassword) (string, error) {
	url := fmt.Sprintf("%s/passwordvault/api/Accounts/%s/Secret/Retrieve", c.BaseURL, accountID)

	response, err := httpJson.SendRequestRaw(url, "POST", c.SessionToken, request, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return "", fmt.Errorf("Failed to retrieve the account SSH Key '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return strings.Trim(string(response), "\""), nil
}

// VerifyAccountCredentials marks an account for verification
func (c Client) VerifyAccountCredentials(accountID string) error {
	url := fmt.Sprintf("%s/passwordvault/API/Accounts/%s/Verify", c.BaseURL, accountID)
	response, err := httpJson.Post(url, c.SessionToken, nil, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to verify account '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return nil
}

// ChangeAccountCredentials marks an account for immediate change
func (c Client) ChangeAccountCredentials(accountID string, changeEntireGroup bool) error {
	url := fmt.Sprintf("%s/passwordvault/API/Accounts/%s/Change", c.BaseURL, accountID)
	body := requests.ChangeAccountCredential{
		ChangeEntireGroup: changeEntireGroup,
	}
	response, err := httpJson.Post(url, c.SessionToken, body, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to mark change on account '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return nil
}

// ReconileAccountCredentials marks an account for reconciliation
func (c Client) ReconileAccountCredentials(accountID string) error {
	url := fmt.Sprintf("%s/passwordvault/API/Accounts/%s/Reconcile", c.BaseURL, accountID)
	response, err := httpJson.Post(url, c.SessionToken, nil, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to mark reconcile on account '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return nil
}
