package ispss

import (
	"context"
	"encoding/json"
	"fmt"

	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/util"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/ispss/responses"
)

// PlatformDiscovery uses the ISPSS API to discover the platform URLs
func PlatformDiscovery(platformURL string) (*responses.PlatformDiscovery, error) {
	ctx := context.TODO()

	subdomain, err := util.GetSubDomain(platformURL)
	if err != nil {
		return &responses.PlatformDiscovery{}, fmt.Errorf("Failed to get subdomain. %s", err)
	}

	url := fmt.Sprintf("https://platform-discovery.cyberark.cloud/api/v2/services/subdomain/%s", subdomain)
	fmt.Printf("URL: %s\n", url)
	fmt.Println("Sending request...")
	_, response, err := httpJson.SendRequestRaw(ctx, false, url, "GET", "", nil, false, nil)
	fmt.Printf("Response: %s\n", response)
	fmt.Printf("Error: %s\n", err)
	if err != nil {
		return &responses.PlatformDiscovery{}, fmt.Errorf("Failed to get platform discovery. %s", err)
	}

	PlatformDiscoveryResponse := &responses.PlatformDiscovery{}
	err = json.Unmarshal(response, PlatformDiscoveryResponse)

	return PlatformDiscoveryResponse, err
}
