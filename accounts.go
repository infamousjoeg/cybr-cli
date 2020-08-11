package pasapi

import (
	"fmt"
	"net/http"
)

// AccountsList contains the count and Accounts array
type AccountsList struct {
	Count    int       `json:"count"`
	Accounts []Account `json:"accounts"`
}

// Account contains the object data for each account returned
type Account struct {
	CategoryModificationTime int                    `json:"categoryModificationTime"`
	ID                       string                 `json:"id"`
	Name                     string                 `json:"name"`
	Address                  string                 `json:"address"`
	Username                 string                 `json:"username"`
	PlatformID               string                 `json:"platformId"`
	SafeName                 string                 `json:"safeName"`
	SecretType               string                 `json:"secretType"`
	PlatformAccountProps     []PlatformAccountProps `json:"platformAccountProps"`
	SecretManagement         []SecretManagement     `json:"secretManagement"`
	CreatedTime              int                    `json:"createdTime"`
}

// PlatformAccountProps contains all additional platform properties data
type PlatformAccountProps struct {
	LogonDomain       string `json:"logonDomain"`
	AWSAccessKeyID    string `json:"awsAccessKey"`
	AWSAccountID      string `json:"awsAccount"`
	AWSARNRole        string `json:"awsARNRole"`
	Port              int    `json:"port"`
	Database          string `json:"database"`
	Index             string `json:"index"`
	DualAccountStatus string `json:"dualAccountStatus"`
	VirtualUsername   string `json:"virtualUsername"`
}

// SecretManagement contains data on the management of the account object returned
type SecretManagement struct {
	AutomaticManagementEnabled bool   `json:"automaticManagementEnabled"`
	ManualManagementReason     string `json:"manualManagementReason"`
	LastModifiedTime           int    `json:"lastModifiedTime"`
}

// AccountsOptions contains all the URL parameter data for the request
type AccountsOptions struct {
	Search     string `json:"search"`
	SearchType string `json:"searchType"`
	Sort       string `json:"sort"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
	Filter     string `json:"filter"`
}

// GetAccounts sends a GET request to the /api/accounts endpoint
func (c *Client) GetAccounts(options *AccountsOptions) (*AccountsList, error) {
	var search string
	var searchType string
	var sort string
	var offset int
	var limit int
	var filter string

	if options != nil {
		search = options.Search
		searchType = options.SearchType
		sort = options.Sort
		offset = options.Offset
		limit = options.Limit
		filter = options.Filter
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accounts?search=%s&searchType=%s&sort=%s&offset=%d&limit=%d&filter=%s", c.BaseURL, search, searchType, sort, offset, limit, filter), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json charset=utf-8")

	res := AccountsList{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
