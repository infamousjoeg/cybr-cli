package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// ListSafesResponse contains an array of all safes the current user can read
type ListSafesResponse struct {
	Safes []ListSafe `json:"Safes"`
}

// ListSafe contains the safe details of every safe the current user can read
// for ListSafesResponse struct
type ListSafe struct {
	SafeURLId   string `json:"SafeUrlId"`
	SafeName    string `json:"SafeName"`
	Description string `json:"Description,omitempty"`
	Location    string `json:"Location"`
}

// ListSafes CyberArk user has access to
func (c Client) ListSafes() (*ListSafesResponse, error) {
	url := fmt.Sprintf("%s/PasswordVault/api/safes", c.BaseURL)
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &ListSafesResponse{}, fmt.Errorf("Failed to list safes. %s", err)
	}
	jsonString, _ := json.Marshal(response)
	ListSafesResponse := ListSafesResponse{}
	err = json.Unmarshal(jsonString, &ListSafesResponse)
	return &ListSafesResponse, err
}

// ListSafeMembersResponse contains data of all members of a specific safe
type ListSafeMembersResponse struct {
	Members []Members `json:"members"`
}

// Members contains all safe member username/group name and their permissions
type Members struct {
	Permissions Permissions `json:"Permissions"`
	Username    string      `json:"UserName"`
}

// Permissions contains the permissions of each safe member
type Permissions struct {
	Add                 bool `json:"Add"`
	AddRenameFolder     bool `json:"AddRenameFolder"`
	BackupSafe          bool `json:"BackupSafe"`
	Delete              bool `json:"Delete"`
	DeleteFolder        bool `json:"DeleteFolder"`
	ListContent         bool `json:"ListContent"`
	ManageSafe          bool `json:"ManageSafe"`
	ManageSafeMembers   bool `json:"ManageSafeMembers"`
	MoveFilesAndFolders bool `json:"MoveFilesAndFolders"`
	Rename              bool `json:"Rename"`
	RestrictedRetrieve  bool `json:"RestrictedRetrieve"`
	Retrieve            bool `json:"Retrieve"`
	Unlock              bool `json:"Unlock"`
	Update              bool `json:"Update"`
	UpdateMetadata      bool `json:"UpdateMetadata"`
	ValidateSafeContent bool `json:"ValidateSafeContent"`
	ViewAudit           bool `json:"ViewAudit"`
	ViewMembers         bool `json:"ViewMembers"`
}

// ListSafeMembers List all members of a safe
func (c Client) ListSafeMembers(safeName string) (*ListSafeMembersResponse, error) {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Safes/%s/Members", c.BaseURL, url.QueryEscape(safeName))
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &ListSafeMembersResponse{}, fmt.Errorf("Failed to list members of safe '%s'. %s", safeName, err)
	}
	jsonString, _ := json.Marshal(response)
	ListSafeMembersResponse := ListSafeMembersResponse{}
	err = json.Unmarshal(jsonString, &ListSafeMembersResponse)
	return &ListSafeMembersResponse, err
}

// AddSafeMemberRequest request sent for adding a member to safe with specific permissions
type AddSafeMemberRequest struct {
	Member AddSafeMember `json:"member"`
}

// AddSafeMember used in AddSafeMemberRequest
type AddSafeMember struct {
	MemberName               string               `json:"MemberName"`
	SearchIn                 string               `json:"SearchIn"`
	MembershipExpirationDate string               `json:"MembershipExpirationDate,omitempty"`
	Permissions              []PermissionKeyValue `json:"Permissions,omitempty"`
}

// PermissionKeyValue used in AddSafeMember struct
type PermissionKeyValue struct {
	Key   string `json:"Key"`
	Value bool   `json:"Value"`
}

// AddSafeMember Add a user or application as a member to a safe with specific permissions
func (c Client) AddSafeMember(safeName string, addMember AddSafeMemberRequest) error {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Safes/%s/Members", c.BaseURL, url.QueryEscape(safeName))
	response, err := httpJson.Post(url, c.SessionToken, addMember, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to add member '%s' to safe '%s'. %s. %s", addMember.Member.MemberName, safeName, string(returnedError), err)
	}
	return nil
}

// RemoveSafeMember Remove a member from a specific safe
func (c Client) RemoveSafeMember(safeName string, member string) error {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Safes/%s/Members/%s", c.BaseURL, url.QueryEscape(safeName), url.QueryEscape(member))
	response, err := httpJson.Delete(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to remove member '%s' from safe '%s'. %s.  %s", member, safeName, string(returnedError), err)
	}
	return nil
}

// AddSafeRequest contains the body of the Add Safe function's request
type AddSafeRequest struct {
	SafeName              string `json:"SafeName"`
	Description           string `json:"Description"`
	OLACEnabled           bool   `json:"OLACEnabled,omitempty"`
	ManagingCPM           string `json:"ManagingCPM"`
	NumberOfDaysRetention int    `json:"NumberOfDaysRetention"`
	AutoPurgeEnabled      bool   `json:"AutoPurgeEnabled,omitempty"`
	SafeLocation          string `json:"Location,omitempty"`
}

// AddSafe to Secure Digital Vault via PAS REST API
func (c Client) AddSafe(body AddSafeRequest) error {
	// Set URL for request
	url := fmt.Sprintf("%s/PasswordVault/api/safes", c.BaseURL)
	_, err := httpJson.Post(url, c.SessionToken, body, c.InsecureTLS, c.Logger)
	if err != nil {
		return fmt.Errorf("Unable to add the safe named %s. %s", body.SafeName, err)
	}
	return nil
}

// DeleteSafe will remove the safeName given to the function via PAS REST API
func (c Client) DeleteSafe(safeName string) error {
	// Set URL for request
	url := fmt.Sprintf("%s/PasswordVault/api/safes/%s", c.BaseURL, url.QueryEscape(safeName))
	_, err := httpJson.Delete(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return fmt.Errorf("Unable to delete the safe named %s. %s", safeName, err)
	}
	return nil
}

// UpdateSafeRequest contains the body of the Update Safe function's request
type UpdateSafeRequest struct {
	SafeName    string `json:"SafeName,omitempty"`
	Description string `json:"Description,omitempty"`
	OLACEnabled bool   `json:"OLACEnabled,omitempty"`
	ManagingCPM string `json:"ManagingCPM,omitempty"`
}

// UpdateSafeResponse contains the response to the Update Safe function's request
type UpdateSafeResponse struct {
	SafeName                  string `json:"SafeName"`
	Description               string `json:"Description"`
	NumberOfDaysRetention     int    `json:"NumberOfDaysRetention"`
	NumberOfVersionsRetention int    `json:"NumberOfVersionsRetention"`
	OLACEnabled               bool   `json:"OLACEnabled"`
}

// UpdateSafe will update the safe's properties that are given for modification
func (c Client) UpdateSafe(targetSafeName string, body UpdateSafeRequest) (*UpdateSafeResponse, error) {
	// Set URL for request
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Safes/%s", c.BaseURL, targetSafeName)
	response, err := httpJson.Put(url, c.SessionToken, body, c.InsecureTLS, c.Logger)
	if err != nil {
		return nil, fmt.Errorf("Unable to update the safe named %s. %s", targetSafeName, err)
	}
	jsonString, _ := json.Marshal(response)
	UpdateSafeResponse := UpdateSafeResponse{}
	err = json.Unmarshal(jsonString, &UpdateSafeResponse)
	return &UpdateSafeResponse, nil
}
