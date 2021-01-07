package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
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
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications?Location=%s", c.BaseURL, location)
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS)
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
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications/%s/Authentications", c.BaseURL, appID)
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS)
	if err != nil {
		return &ListApplicationAuthenticationMethodsResponse{}, fmt.Errorf("Error listing application's '%s' authentication methods. %s", appID, err)
	}
	jsonString, _ := json.Marshal(response)
	ListApplicationAuthenticationMethodsResponse := ListApplicationAuthenticationMethodsResponse{}
	err = json.Unmarshal(jsonString, &ListApplicationAuthenticationMethodsResponse)
	return &ListApplicationAuthenticationMethodsResponse, err
}

// AddApplicationRequest request for adding application
type AddApplicationRequest struct {
	Application Application `json:"application"`
}

// Application is used when adding an application
type Application struct {
	AppID               string `json:"AppID"`
	Description         string `json:"Description,omitempty"`
	Location            string `json:"Location,omitempty"`
	AccessPermittedFrom int    `json:"AccessPermittedFrom"`
	AccessPermittedTo   int    `json:"AccessPermittedTo"`
	ExpirationDate      string `json:"ExpirationDate,omitempty"`
	Disabled            string `json:"Disabled,omitempty"`
	BusinessOwnerFName  string `json:"BusinessOwnerFName,omitempty"`
	BusinessOwnerLName  string `json:"BusinessOwnerLName,omitempty"`
	BusinessOwnerEmail  string `json:"BusinessOwnerEmail,omitempty"`
	BusinessOwnerPhone  string `json:"BusinessOwnerPhone,omitempty"`
}

// AddApplication add an applications to PAS
func (c Client) AddApplication(application AddApplicationRequest) error {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications", c.BaseURL)
	response, err := httpJson.Post(url, c.SessionToken, application, c.InsecureTLS)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Error adding application '%s' authentication methods. %s. %s", application.Application.AppID, string(returnedError), err)
	}
	return nil
}

// DeleteApplication delete an applications to PAS
func (c Client) DeleteApplication(appID string) error {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications/%s", c.BaseURL, url.QueryEscape(appID))
	response, err := httpJson.Delete(url, c.SessionToken, c.InsecureTLS)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Error deleting application '%s' authentication methods. %s. %s", appID, string(returnedError), err)
	}
	return nil
}

// AddApplicationAuthenticationRequest request to add authentication method to application
type AddApplicationAuthenticationRequest struct {
	Authentication ApplicationAuthenticationMethod `json:"authentication"`
}

// ApplicationAuthenticationMethod represents an application authentication method
type ApplicationAuthenticationMethod struct {
	AuthType             string `json:"AuthType"`
	AuthValue            string `json:"AuthValue"`
	IsFolder             bool   `json:"IsFolder,omitempty"`
	AllowInternalScripts bool   `json:"AllowInternalScripts,omitempty"`
}

// AddApplicationAuthenticationMethod add authentication method to an application
func (c Client) AddApplicationAuthenticationMethod(appID string, authenticationMethod AddApplicationAuthenticationRequest) error {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications/%s/Authentications/", c.BaseURL, url.QueryEscape(appID))
	response, err := httpJson.Post(url, c.SessionToken, authenticationMethod, c.InsecureTLS)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Error adding application authentication method to '%s'. %s. %s", appID, string(returnedError), err)
	}
	return nil
}

// DeleteApplicationAuthenticationMethod delete an applications authentication method
func (c Client) DeleteApplicationAuthenticationMethod(appID string, authnMethodID string) error {
	url := fmt.Sprintf("%s/PasswordVault/WebServices/PIMServices.svc/Applications/%s/Authentications/%s", c.BaseURL, url.QueryEscape(appID), url.QueryEscape(authnMethodID))
	response, err := httpJson.Delete(url, c.SessionToken, c.InsecureTLS)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Error deleting application '%s' authentication methods. %s. %s", appID, string(returnedError), err)
	}
	return nil
}
