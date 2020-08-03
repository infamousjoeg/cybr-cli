package pasapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// AccountsList contains the count and Accounts array
type AccountsList struct {
	Count    int        `json:"count"`
	Accounts []Accounts `json:"accounts"`
}

// Accounts contains the object data for each account returned
type Accounts struct {
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
func (c *Client) GetAccounts(ctx context.Context, options *AccountsOptions) (*AccountsList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/accounts", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := AccountsList{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", c.sessionToken)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	fullResponse := successResponse{
		Data: v,
	}

	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}

	return nil
}
