package conjur

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Whoami gets current user info logged in to Conjur
func Whoami() (map[string]interface{}, error) {
	client, loginPair, err := GetConjurClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize conjur client. %s", err)
	}

	config := client.GetConfig()
	url := fmt.Sprintf("%s/whoami", config.ApplianceURL)
	resp, err := sendConjurAuthenticatedHTTPRequest(client, loginPair, url, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to get logged in user. %s", err)
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read conjur whoami body. %s", err)
	}

	resultInterface := make(map[string]interface{})
	err = json.Unmarshal(result, &resultInterface)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal conjur whoami body into json object. %s", err)
	}

	return resultInterface, nil
}
