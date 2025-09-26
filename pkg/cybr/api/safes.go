package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/queries"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/responses"
	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// GetSafeDetails get safe details and optionally include accounts
func (c Client) GetSafeDetails(safeName string, includeAccounts bool) (*responses.GetSafeDetails, error) {
	if len(safeName) == 0 {
		return &responses.GetSafeDetails{}, fmt.Errorf("no safe name passed in")
	}
	// https://<subdomain>.privilegecloud.cyberark.cloud/PasswordVault/API/Safes/{SafeUrlId}/
	url := fmt.Sprintf("%s/passwordvault/api/Safes/%s/", c.BaseURL, url.QueryEscape(safeName))
	if includeAccounts {
		url = fmt.Sprintf("%s?includeAccounts=true", url)
	}
	response, err := httpJson.Get(false, url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.GetSafeDetails{}, fmt.Errorf("failed to get safe details. %s", err)
	}
	jsonString, _ := json.Marshal(response)
	c.Logger.Writef("GET %s\n%s\n\n", url, jsonString)

	GetSafeDetailsResponse := responses.GetSafeDetails{}
	err = json.Unmarshal(jsonString, &GetSafeDetailsResponse)
	return &GetSafeDetailsResponse, err
}

// ListSafes CyberArk user has access to
func (c Client) ListSafes() (*responses.ListSafes, error) {
	url := fmt.Sprintf("%s/passwordvault/api/safes", c.BaseURL)
	response, err := httpJson.Get(false, url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.ListSafes{}, fmt.Errorf("failed to list safes. %s", err)
	}
	jsonString, _ := json.Marshal(response)
	ListSafesResponse := responses.ListSafes{}
	err = json.Unmarshal(jsonString, &ListSafesResponse)
	return &ListSafesResponse, err
}

// ListSafeMembers List all members of a safe
func (c Client) ListSafeMembers(safeName string, query *queries.ListSafeMembers) (*responses.ListSafeMembers, error) {
	url := fmt.Sprintf("%s/passwordvault/api/Safes/%s/Members%s", c.BaseURL, url.QueryEscape(safeName), httpJson.GetURLQuery(query))
	response, err := httpJson.Get(false, url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.ListSafeMembers{}, fmt.Errorf("failed to list members of safe '%s'. %s", safeName, err)
	}
	jsonString, _ := json.Marshal(response)
	ListSafeMembersResponse := responses.ListSafeMembers{}
	err = json.Unmarshal(jsonString, &ListSafeMembersResponse)
	return &ListSafeMembersResponse, err
}

// AddSafeMember Add a user or application as a member to a safe with specific permissions
func (c Client) AddSafeMember(safeName string, addMember requests.AddSafeMember) error {
	url := fmt.Sprintf("%s/passwordvault/api/safes/%s/members", c.BaseURL, url.QueryEscape(safeName))
	response, err := httpJson.Post(false, url, c.SessionToken, addMember, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("failed to add member '%s' to safe '%s'. %s. %s", addMember.MemberName, safeName, string(returnedError), err)
	}
	return nil
}

// RemoveSafeMember Remove a member from a specific safe
func (c Client) RemoveSafeMember(safeName string, member string) error {
	url := fmt.Sprintf("%s/passwordvault/api/Safes/%s/Members/%s", c.BaseURL, url.QueryEscape(safeName), url.QueryEscape(member))
	response, err := httpJson.Delete(false, url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("failed to remove member '%s' from safe '%s'. %s.  %s", member, safeName, string(returnedError), err)
	}
	return nil
}

// UpdateSafeMember Update safe member permissions
func (c Client) UpdateSafeMember(safeName string, memberName string, updateMember requests.UpdateSafeMember) error {
	// PUT https://<subdomain>.privilegecloud.cyberark.cloud/PasswordVault/API/Safes/{SafeUrlId}/Members/{MemberName}/
	url := fmt.Sprintf("%s/passwordvault/api/safes/%s/members/%s/", c.BaseURL, url.QueryEscape(safeName), url.QueryEscape(memberName))
	jsonString, _ := json.Marshal(updateMember)
	c.Logger.Writef("Update Request:\n%s\n", jsonString)

	response, err := httpJson.Put(false, url, c.SessionToken, updateMember, c.InsecureTLS, c.Logger)
	jsonString, _ = json.Marshal(response)
	c.Logger.Writef("PUT %s\nResponse:\n%s\n\n", url, jsonString)

	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("failed to update member '%s' to safe '%s'. %s. %s", memberName, safeName, string(returnedError), err)
	}
	return nil
}

// AddSafe to Secure Digital Vault via PAS REST API
func (c Client) AddSafe(body requests.AddSafe) error {
	// Set URL for request
	url := fmt.Sprintf("%s/passwordvault/api/safes", c.BaseURL)
	_, err := httpJson.Post(false, url, c.SessionToken, body, c.InsecureTLS, c.Logger)
	if err != nil {
		return fmt.Errorf("unable to add the safe named %s. %s", body.SafeName, err)
	}
	return nil
}

// DeleteSafe will remove the safeName given to the function via PAS REST API
func (c Client) DeleteSafe(safeName string) error {
	// Set URL for request
	url := fmt.Sprintf("%s/passwordvault/api/safes/%s", c.BaseURL, url.QueryEscape(safeName))
	_, err := httpJson.Delete(false, url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return fmt.Errorf("unable to delete the safe named %s. %s", safeName, err)
	}
	return nil
}

// UpdateSafe will update the safe's properties that are given for modification
func (c Client) UpdateSafe(targetSafeName string, body requests.UpdateSafe) (*responses.UpdateSafe, error) {
	// Set URL for request
	url := fmt.Sprintf("%s/passwordvault/WebServices/PIMServices.svc/Safes/%s", c.BaseURL, targetSafeName)
	response, err := httpJson.Put(false, url, c.SessionToken, body, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.UpdateSafe{}, fmt.Errorf("unable to update the safe named %s. %s", targetSafeName, err)
	}
	jsonString, _ := json.Marshal(response)
	UpdateSafeResponse := responses.UpdateSafe{}
	err = json.Unmarshal(jsonString, &UpdateSafeResponse)
	if err != nil {
		return &responses.UpdateSafe{}, fmt.Errorf("failed to parse response for update safe named %s, : %s", targetSafeName, err)
	}
	return &UpdateSafeResponse, nil
}

// FilterSafes will return a list of safes that match the given filter, commonly used to filter by safe member
func (c Client) FilterSafes(filter string, search string) ([]string, error) {
	// List All Safes
	safes, err := c.ListSafes()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a list of all safes. %s", err)
	}

	// For each safe, extract the safe name and ListSafeMembers for that safe
	query := &queries.ListSafeMembers{
		Filter: filter,
		Search: search,
	}

	var filteredSafes []string
	for i := 0; i < len(safes.Safes); i++ {
		safeName := safes.Safes[i].SafeName
		if safeName == "Notification Engine" {
			continue
		}
		// List Safe Members
		listMemberResult, err := c.ListSafeMembers(safeName, query)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve a list of members for safe '%s'. %s", safeName, err)
		}
		// Unmarshal the response
		jsonListMemberResult, _ := json.Marshal(listMemberResult)
		parsedListMemberResult := responses.ListSafeMembers{}
		err = json.Unmarshal(jsonListMemberResult, &parsedListMemberResult)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal the safe member from %s. %s", safeName, err)
		}
		if parsedListMemberResult.Count > 0 {
			filteredSafes = append(filteredSafes, safes.Safes[i].SafeName)
		}
	}

	return filteredSafes, nil
}
