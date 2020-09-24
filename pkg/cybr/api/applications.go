package api

import (
	"encoding/json"
	"fmt"
	"time"

	httpJson "github.com/infamousjoeg/pas-api-go/pkg/cybr/helpers"
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
func (c Client) ListApplications(location string) (*ListApplicationsResponse, error) {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications?Location=%s", c.Hostname, location)
	response, err := httpJson.Get(url, c.sessionToken)
	if err != nil {
		return &ListApplicationsResponse{}, fmt.Errorf("Error listing applications in location '%s'. %s", location, err)
	}
	jsonString, _ := json.Marshal(response)
	ListApplicationsResponse := ListApplicationsResponse{}
	err = json.Unmarshal(jsonString, &ListApplicationsResponse)
	return &ListApplicationsResponse, err
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
	IsFolder             bool   `json:"IsFolder,omitempty"`
	AuthID               string `json:"authID"`
}

// ListApplicationAuthenticationMethods returns all auth methods for a specific Application Identity
func (c Client) ListApplicationAuthenticationMethods(appID string) (*ListApplicationAuthenticationMethodsResponse, error) {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications/%s/Authentications", c.Hostname, appID)
	response, err := httpJson.Get(url, c.sessionToken)
	if err != nil {
		return &ListApplicationAuthenticationMethodsResponse{}, fmt.Errorf("Error listing application's '%s' authentication methods. %s", appID, err)
	}
	jsonString, _ := json.Marshal(response)
	ListApplicationAuthenticationMethodsResponse := ListApplicationAuthenticationMethodsResponse{}
	err = json.Unmarshal(jsonString, &ListApplicationAuthenticationMethodsResponse)
	return &ListApplicationAuthenticationMethodsResponse, err
}
