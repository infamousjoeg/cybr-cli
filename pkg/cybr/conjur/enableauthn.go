package conjur

import (
	"fmt"
	"net/url"
	"strings"
)

// EnableAuthenticator enable a specific authenticator in conjur. e.g. serviceID: authn-iam/prod or authn-k8s/k8s-cluster-1
func EnableAuthenticator(serviceID string) error {
	client, loginPair, err := GetConjurClient()
	if err != nil {
		return fmt.Errorf("Failed to initialize conjur client. %s", err)
	}

	config := client.GetConfig()
	url := fmt.Sprintf("%s/%s/%s", config.ApplianceURL, serviceID, url.QueryEscape(config.Account))
	body := strings.NewReader("enabled=true")
	_, err = sendConjurAuthenticatedHTTPRequest(client, loginPair, url, "PATCH", body)
	if err != nil {
		return fmt.Errorf("Failed to enable authenticator '%s'. %s", serviceID, err)
	}
	return nil
}
