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

// ListSafes CyberArk user has access to
func (c Client) ListSafes() (*responses.ListSafes, error) {
	url := fmt.Sprintf("%s/passwordvault/api/safes", c.BaseURL)
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.ListSafes{}, fmt.Errorf("Failed to list safes. %s", err)
	}
	jsonString, _ := json.Marshal(response)
	ListSafesResponse := responses.ListSafes{}
	err = json.Unmarshal(jsonString, &ListSafesResponse)
	return &ListSafesResponse, err
}

// ListSafeMembers List all members of a safe
func (c Client) ListSafeMembers(safeName string, query *queries.ListSafeMembers) (*responses.ListSafeMembers, error) {
	url := fmt.Sprintf("%s/passwordvault/api/Safes/%s/Members%s", c.BaseURL, url.QueryEscape(safeName), httpJson.GetURLQuery(query))
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.ListSafeMembers{}, fmt.Errorf("Failed to list members of safe '%s'. %s", safeName, err)
	}
	jsonString, _ := json.Marshal(response)
	ListSafeMembersResponse := responses.ListSafeMembers{}
	err = json.Unmarshal(jsonString, &ListSafeMembersResponse)
	return &ListSafeMembersResponse, err
}

// AddSafeMember Add a user or application as a member to a safe with specific permissions
func (c Client) AddSafeMember(safeName string, addMember requests.AddSafeMember) error {
	url := fmt.Sprintf("%s/passwordvault/api/safes/%s/members", c.BaseURL, url.QueryEscape(safeName))
	response, err := httpJson.Post(url, c.SessionToken, addMember, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to add member '%s' to safe '%s'. %s. %s", addMember.MemberName, safeName, string(returnedError), err)
	}
	return nil
}

// RemoveSafeMember Remove a member from a specific safe
func (c Client) RemoveSafeMember(safeName string, member string) error {
	url := fmt.Sprintf("%s/passwordvault/WebServices/PIMServices.svc/Safes/%s/Members/%s", c.BaseURL, url.QueryEscape(safeName), url.QueryEscape(member))
	response, err := httpJson.Delete(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		returnedError, _ := json.Marshal(response)
		return fmt.Errorf("Failed to remove member '%s' from safe '%s'. %s.  %s", member, safeName, string(returnedError), err)
	}
	return nil
}

// AddSafe to Secure Digital Vault via PAS REST API
func (c Client) AddSafe(body requests.AddSafe) error {
	// Set URL for request
	url := fmt.Sprintf("%s/passwordvault/api/safes", c.BaseURL)
	_, err := httpJson.Post(url, c.SessionToken, body, c.InsecureTLS, c.Logger)
	if err != nil {
		return fmt.Errorf("Unable to add the safe named %s. %s", body.SafeName, err)
	}
	return nil
}

// DeleteSafe will remove the safeName given to the function via PAS REST API
func (c Client) DeleteSafe(safeName string) error {
	// Set URL for request
	url := fmt.Sprintf("%s/passwordvault/api/safes/%s", c.BaseURL, url.QueryEscape(safeName))
	_, err := httpJson.Delete(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return fmt.Errorf("Unable to delete the safe named %s. %s", safeName, err)
	}
	return nil
}

// UpdateSafe will update the safe's properties that are given for modification
func (c Client) UpdateSafe(targetSafeName string, body requests.UpdateSafe) (*responses.UpdateSafe, error) {
	// Set URL for request
	url := fmt.Sprintf("%s/passwordvault/WebServices/PIMServices.svc/Safes/%s", c.BaseURL, targetSafeName)
	response, err := httpJson.Put(url, c.SessionToken, body, c.InsecureTLS, c.Logger)
	if err != nil {
		return nil, fmt.Errorf("Unable to update the safe named %s. %s", targetSafeName, err)
	}
	jsonString, _ := json.Marshal(response)
	UpdateSafeResponse := responses.UpdateSafe{}
	err = json.Unmarshal(jsonString, &UpdateSafeResponse)
	return &UpdateSafeResponse, nil
}

// FilterSafes will return a list of safes that match the given filter, commonly used to filter by safe member
func (c Client) FilterSafes(filter string, search string) ([]string, error) {
	// List All Safes
	safes, err := c.ListSafes()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve a list of all safes. %s", err)
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
			return nil, fmt.Errorf("Failed to retrieve a list of members for safe '%s'. %s", safeName, err)
		}
		// Unmarshal the response
		jsonListMemberResult, _ := json.Marshal(listMemberResult)
		parsedListMemberResult := responses.ListSafeMembers{}
		err = json.Unmarshal(jsonListMemberResult, &parsedListMemberResult)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal the safe member from %s. %s", safeName, err)
		}
		if parsedListMemberResult.Count > 0 {
			filteredSafes = append(filteredSafes, safes.Safes[i].SafeName)
		}
	}

	return filteredSafes, nil
}
