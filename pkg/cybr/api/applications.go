package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/responses"
	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// ListApplications returns all Application Identities setup in PAS
func (c Client) ListApplications(location string) (*responses.ListApplications, error) {
	url := fmt.Sprintf("%s/passwordvault/WebServices/PIMServices.svc/Applications?Location=%s", c.BaseURL, location)
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.ListApplications{}, fmt.Errorf("Error listing applications in location '%s'. %s", location, err)
	}
	jsonString, _ := json.Marshal(response)
	ListApplicationsResponse := responses.ListApplications{}
	err = json.Unmarshal(jsonString, &ListApplicationsResponse)
	return &ListApplicationsResponse, err
}

// ListApplicationAuthenticationMethods returns all auth methods for a specific Application Identity
func (c Client) ListApplicationAuthenticationMethods(appID string) (*responses.ListApplicationAuthenticationMethods, error) {
	url := fmt.Sprintf("%s/passwordvault/WebServices/PIMServices.svc/Applications/%s/Authentications", c.BaseURL, appID)
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.ListApplicationAuthenticationMethods{}, fmt.Errorf("Error listing application's '%s' authentication methods. %s", appID, err)
	}
	jsonString, _ := json.Marshal(response)
	ListApplicationAuthenticationMethodsResponse := responses.ListApplicationAuthenticationMethods{}
	err = json.Unmarshal(jsonString, &ListApplicationAuthenticationMethodsResponse)
	return &ListApplicationAuthenticationMethodsResponse, err
}

// AddApplication add an applications to PAS
func (c Client) AddApplication(application requests.AddApplication) error {
	url := fmt.Sprintf("%s/passwordvault/WebServices/PIMServices.svc/Applications", c.BaseURL)
	response, err := httpJson.Post(url, c.SessionToken, application, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Error adding application '%s' authentication methods. %s. %s", application.Application.AppID, string(returnedError), err)
	}
	return nil
}

// DeleteApplication delete an applications to PAS
func (c Client) DeleteApplication(appID string) error {
	url := fmt.Sprintf("%s/passwordvault/WebServices/PIMServices.svc/Applications/%s", c.BaseURL, url.QueryEscape(appID))
	response, err := httpJson.Delete(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Error deleting application '%s' authentication methods. %s. %s", appID, string(returnedError), err)
	}
	return nil
}

// AddApplicationAuthenticationMethod add authentication method to an application
func (c Client) AddApplicationAuthenticationMethod(appID string, authenticationMethod requests.AddApplicationAuthentication) error {
	url := fmt.Sprintf("%s/passwordvault/WebServices/PIMServices.svc/Applications/%s/Authentications/", c.BaseURL, url.QueryEscape(appID))
	response, err := httpJson.Post(url, c.SessionToken, authenticationMethod, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Error adding application authentication method to '%s'. %s. %s", appID, string(returnedError), err)
	}
	return nil
}

// DeleteApplicationAuthenticationMethod delete an applications authentication method
func (c Client) DeleteApplicationAuthenticationMethod(appID string, authnMethodID string) error {
	url := fmt.Sprintf("%s/passwordvault/WebServices/PIMServices.svc/Applications/%s/Authentications/%s", c.BaseURL, url.QueryEscape(appID), url.QueryEscape(authnMethodID))
	response, err := httpJson.Delete(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Error deleting application '%s' authentication methods. %s. %s", appID, string(returnedError), err)
	}
	return nil
}
