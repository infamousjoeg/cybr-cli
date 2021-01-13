package conjur

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Info get info of the conjur instance
func Info() (map[string]interface{}, error) {
	client, loginPair, err := GetConjurClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize conjur client. %s", err)
	}

	config := client.GetConfig()
	url := fmt.Sprintf("%s/info", config.ApplianceURL)
	resp, err := sendConjurAuthenticatedHTTPRequest(client, loginPair, url, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to get info. %s", err)
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read conjur info body. %s", err)
	}

	resultInterface := make(map[string]interface{})
	err = json.Unmarshal(result, &resultInterface)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal conjur info body into json object. %s", err)
	}

	return resultInterface, nil
}
