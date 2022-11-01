package api

import (
	"encoding/json"
	"fmt"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/queries"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/responses"
	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// ListPlatforms available in CyberArk
func (c Client) ListPlatforms(query *queries.ListPlatforms) (*responses.ListPlatforms, error) {
	url := fmt.Sprintf("%s/passwordvault/api/platforms%s", c.BaseURL, httpJson.GetURLQuery(query))
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.ListPlatforms{}, fmt.Errorf("Failed to list platforms. %s", err)
	}
	jsonString, _ := json.Marshal(response)
	ListPlatformsResponse := responses.ListPlatforms{}
	err = json.Unmarshal(jsonString, &ListPlatformsResponse)
	return &ListPlatformsResponse, err
}

// GetPlatform details for specific platform
func (c Client) GetPlatform(platformID string) (*responses.GetPlatform, error) {
	url := fmt.Sprintf("%s/passwordvault/api/platforms/%s", c.BaseURL, platformID)
	response, err := httpJson.Get(url, c.SessionToken, c.InsecureTLS, c.Logger)
	if err != nil {
		return &responses.GetPlatform{}, fmt.Errorf("Failed to get platform. %s", err)
	}

	jsonString, _ := json.Marshal(response)
	GetPlatformResponse := &responses.GetPlatform{}
	err = json.Unmarshal(jsonString, GetPlatformResponse)
	return GetPlatformResponse, err
}
