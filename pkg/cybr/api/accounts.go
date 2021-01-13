package api

import (
	"encoding/json"
	"fmt"

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
	response, err := httpJson.Post(url, c.SessionToken, account, c.InsecureTLS, c.Logger)
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
