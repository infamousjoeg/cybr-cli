package api

import (
	"encoding/json"
	"fmt"
	"log"

	httpJson "github.com/infamousjoeg/pas-api-go/pkg/cybr/helpers"
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
func ListSafes(hostname string, token string) *ListSafesResponse {
	url := fmt.Sprintf("%s/PasswordVault/api/safes", hostname)
	response, err := httpJson.Get(url, token)
	if err != nil {
		log.Fatalf("Error listing safes for current CyberArk user. %s", err)
	}
	jsonString, _ := json.Marshal(response)
	ListSafesResponse := ListSafesResponse{}
	json.Unmarshal(jsonString, &ListSafesResponse)
	return &ListSafesResponse
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
func ListSafeMembers(hostname string, token string, safeName string) *ListSafeMembersResponse {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Safes/%s/Members", hostname, safeName)
	response, err := httpJson.Get(url, token)
	if err != nil {
		log.Fatalf("Error listing safe members for %s safe. %s", safeName, err)
	}
	jsonString, _ := json.Marshal(response)
	ListSafeMembersResponse := ListSafeMembersResponse{}
	json.Unmarshal(jsonString, &ListSafeMembersResponse)
	return &ListSafeMembersResponse
}
