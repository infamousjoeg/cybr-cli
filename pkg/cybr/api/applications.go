package api

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	httpJson "github.com/infamousjoeg/pas-api-go/pkg/helpers"
)

//ListApplicationsResponse contains all applications and app data returned
type ListApplicationsResponse struct {
	Application []ListApplication `json:"application"`
}

// ListApplication contains all specific application data for ListApplicationsResponse
type ListApplication struct {
	AccessPermittedFrom                     int       `json:"AccessPermittedFrom"`
	AccessPermittedTo                       int       `json:"AccessPermittedTo"`
	AllowExtendedAuthenticationRestrictions bool      `json:"AllowExtendedAuthenticationRestrict"`
	AppID                                   string    `json:"AppID"`
	BusinessOwnerEmail                      string    `json:"BusinessOwnerEmail"`
	BusinessOwnerFName                      string    `json:"BusinessOwnerFName"`
	BusinessOwnerLName                      string    `json:"BusinessOwnerLName"`
	BusinessOwnerPhone                      string    `json:"BusinessOwnerPhone"`
	Description                             string    `json:"Description"`
	Disabled                                bool      `json:"Disabled"`
	ExpirationDate                          time.Time `json:"ExpirationDate"`
	Location                                string    `json:"Location"`
}

// ListApplications returns all Application Identities setup in PAS
func ListApplications(hostname string, token string, location string, subLocations bool) *ListApplicationsResponse {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications?Location=%s&IncludeSublocations=%t", hostname, location, subLocations)
	response, err := httpJson.Get(url, token)
	if err != nil {
		log.Fatalf("Error listing applications. %s", err)
	}
	jsonString, _ := json.Marshal(response)
	ListApplicationsResponse := ListApplicationsResponse{}
	json.Unmarshal(jsonString, &ListApplicationsResponse)
	return &ListApplicationsResponse
}

// ListApplicationAuthenticationMethodsResponse contains all auth methods for a specific App ID
type ListApplicationAuthenticationMethodsResponse struct {
	Authentication []ListAuthentication `json:"authentication"`
}

// ListAuthentication contains details of the authentication method listed
type ListAuthentication struct {
	AllowInternalScripts bool   `json:"AllowInternalScripts,omitempty"`
	AppID                string `json:"AppID"`
	AuthType             string `json:"AuthType"`
	AuthValue            string `json:"AuthValue"`
	Comment              string `json:"Comment,omitempty"`
	IsFolder             string `json:"IsFolder,omitempty"`
	AuthID               int    `json:"authID"`
}

// ListApplicationAuthenticationMethods returns all auth methods for a specific Application Identity
func ListApplicationAuthenticationMethods(hostname string, token string, appID string) *ListApplicationAuthenticationMethodsResponse {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications/%s/Authentications", hostname, appID)
	response, err := httpJson.Get(url, token)
	if err != nil {
		log.Fatalf("Error listing %s authentication methods. %s", appID, err)
	}
	jsonString, _ := json.Marshal(response)
	ListApplicationAuthenticationMethodsResponse := ListApplicationAuthenticationMethodsResponse{}
	json.Unmarshal(jsonString, &ListApplicationAuthenticationMethodsResponse)
	return &ListApplicationAuthenticationMethodsResponse
}
