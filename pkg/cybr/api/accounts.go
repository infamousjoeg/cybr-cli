package api

import (
	"encoding/json"
	"fmt"
	"strings"

	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// ListAccountQueryParams represents valid query parameters when listing accounts
type ListAccountQueryParams struct {
	Search     string `query_key:"search"`
	SearchType string `query_key:"searchType"`
	Sort       string `query_key:"sort"`
	Offset     int    `query_key:"offset"`
	Limit      int    `query_key:"limit"`
	Filter     string `query_key:"filter"`
}

// ListAccountResponse response from listing accounts
type ListAccountResponse struct {
	Value []GetAccountResponse `json:"value"`
	Count int                  `json:"count"`
}

// GetAccountResponse response from getting specific account details
type GetAccountResponse struct {
	CategoryModificationTime  int                    `json:"categoryModificationTime"`
	ID                        string                 `json:"id"`
	Name                      string                 `json:"name"`
	Address                   string                 `json:"address"`
	UserName                  string                 `json:"userName"`
	PlatformID                string                 `json:"platformId"`
	SafeName                  string                 `json:"safeName"`
	SecretType                string                 `json:"secretType"`
	PlatformAccountProperties map[string]interface{} `json:"platformAccountProperties"`
	SecretManagement          SecretManagement       `json:"secretManagement"`
	CreatedTime               int                    `json:"createdTime"`
}

// SecretManagement used in getting and setting accounts
type SecretManagement struct {
	AutomaticManagementEnabled bool   `json:"automaticManagementEnabled"`
	Status                     string `json:"status"`
	ManualManagementReason     string `json:"manualManagementReason,omitempty"`
	LastModifiedTime           int    `json:"lastModifiedTime,omitempty"`
}

// AddAccountRequest request used to create an account
type AddAccountRequest struct {
	Name                      string            `json:"name,omitempty"`
	Address                   string            `json:"address"`
	UserName                  string            `json:"userName"`
	PlatformID                string            `json:"platformId"`
	SafeName                  string            `json:"safeName"`
	SecretType                string            `json:"secretType"`
	Secret                    string            `json:"secret"`
	PlatformAccountProperties map[string]string `json:"platformAccountProperties,omitempty"`
	SecretManagement          SecretManagement  `json:"secretManagement,omitempty"`
}

// GetAccountPasswordRequest get an accounts password from PAS. All fields are optional expect if ticketing system is integrated
type GetAccountPasswordRequest struct {
	Reason              string `json:"reason,omitempty"`
	TicketingSystemName string `json:"TicketingSystemName,omitempty"`
	TicketID            string `json:"TicketId,omitempty"`
	Version             int    `json:"Version,omitempty"`
	ActionType          string `json:"ActionType,omitempty"`
	IsUse               bool   `json:"isUse,omitempty"`
	Machine             string `json:"Machine,omitempty"`
}

// ChangeAccountCredentialRequest only used when account is part of a group
type ChangeAccountCredentialRequest struct {
	ChangeEntireGroup bool `json:"ChangeEntireGroup,omitempty"`
}

// ListAccounts CyberArk user has access to
func (c Client) ListAccounts(query *ListAccountQueryParams) (*ListAccountResponse, error) {
	url := fmt.Sprintf("%s/PasswordVault/api/Accounts%s", c.BaseURL, httpJson.GetURLQuery(query))
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &ListAccountResponse{}, fmt.Errorf("Failed to list accounts. %s", err)
	}

	jsonString, _ := json.Marshal(response)
	ListSafesResponse := ListAccountResponse{}
	err = json.Unmarshal(jsonString, &ListSafesResponse)
	return &ListSafesResponse, err
}

// GetAccount details for specific account
func (c Client) GetAccount(accountID string) (*GetAccountResponse, error) {
	url := fmt.Sprintf("%s/PasswordVault/api/Accounts/%s", c.BaseURL, accountID)
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &GetAccountResponse{}, fmt.Errorf("Failed to get account. %s", err)
	}

	jsonString, _ := json.Marshal(response)
	GetAccountResponse := &GetAccountResponse{}
	err = json.Unmarshal(jsonString, GetAccountResponse)
	return GetAccountResponse, err
}

// AddAccount to cyberark
func (c Client) AddAccount(account AddAccountRequest) (*GetAccountResponse, error) {
	url := fmt.Sprintf("%s/PasswordVault/api/Accounts", c.BaseURL)
	logger := c.GetLogger().AddSecret(account.Secret)
	response, err := httpJson.Post(url, c.SessionToken, account, c.InsecureTLS, logger)
	logger = logger.ClearSecrets()

	if err != nil {
		returnedError, _ := json.Marshal(response)
		return &GetAccountResponse{}, fmt.Errorf("Failed to add account. %s. %s", string(returnedError), err)
	}

	jsonString, _ := json.Marshal(response)
	GetAccountResponse := &GetAccountResponse{}
	err = json.Unmarshal(jsonString, GetAccountResponse)
	return GetAccountResponse, err
}

// DeleteAccount from cyberark
func (c Client) DeleteAccount(accountID string) error {
	url := fmt.Sprintf("%s/PasswordVault/api/Accounts/%s", c.BaseURL, accountID)
	response, err := httpJson.Delete(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to delete account '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return nil
}

// GetJITAccess from a specific account
func (c Client) GetJITAccess(accountID string) error {
	url := fmt.Sprintf("%s/PasswordVault/api/Accounts/%s/grantAdministrativeAccess", c.BaseURL, accountID)
	response, err := httpJson.Post(url, c.SessionToken, nil, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to get JIT access for account '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return nil
}

// RevokeJITAccess from a specific account
func (c Client) RevokeJITAccess(accountID string) error {
	url := fmt.Sprintf("%s/PasswordVault/api/Accounts/%s/RevokeAdministrativeAccess", c.BaseURL, accountID)
	response, err := httpJson.Post(url, c.SessionToken, nil, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to revoke JIT access for account '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return nil
}

// GetAccountPassword This method enables users to retrieve the password or SSH key of an existing account that is identified by its Account ID. It enables users to specify a reason and ticket ID, if required
func (c Client) GetAccountPassword(accountID string, request GetAccountPasswordRequest) (string, error) {
	url := fmt.Sprintf("%s/PasswordVault/api/Accounts/%s/Password/Retrieve", c.BaseURL, accountID)

	response, err := httpJson.SendRequestRaw(url, "POST", c.SessionToken, request, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return "", fmt.Errorf("Failed to retrieve the account password '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return strings.Trim(string(response), "\""), nil
}

// GetAccountSSHKey This method enables users to retrieve the password or SSH key of an existing account that is identified by its Account ID. It enables users to specify a reason and ticket ID, if required
func (c Client) GetAccountSSHKey(accountID string, request GetAccountPasswordRequest) (string, error) {
	url := fmt.Sprintf("%s/PasswordVault/api/Accounts/%s/Secret/Retrieve", c.BaseURL, accountID)

	response, err := httpJson.SendRequestRaw(url, "POST", c.SessionToken, request, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return "", fmt.Errorf("Failed to retrieve the account SSH Key '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return strings.Trim(string(response), "\""), nil
}

// VerifyAccountCredentials marks an account for verification
func (c Client) VerifyAccountCredentials(accountID string) error {
	url := fmt.Sprintf("%s/PasswordVault/API/Accounts/%s/Verify", c.BaseURL, accountID)
	response, err := httpJson.Post(url, c.SessionToken, nil, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to verify account '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return nil
}

// ChangeAccountCredentials marks an account for immediate change
func (c Client) ChangeAccountCredentials(accountID string, changeEntireGroup bool) error {
	url := fmt.Sprintf("%s/PasswordVault/API/Accounts/%s/Change", c.BaseURL, accountID)
	body := ChangeAccountCredentialRequest{
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
	url := fmt.Sprintf("%s/PasswordVault/API/Accounts/%s/Reconcile", c.BaseURL, accountID)
	response, err := httpJson.Post(url, c.SessionToken, nil, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to mark reconcile on account '%s'. %s. %s", accountID, string(returnedError), err)
	}

	return nil
}
