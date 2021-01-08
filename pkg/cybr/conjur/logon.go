package conjur

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Login to conjur and return an api key in []byte format
func Login(applianceURL string, account string, username string, password []byte, certPath string) ([]byte, error) {
	content, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read cert file. %s. %s", certPath, err)
	}

	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(content)
	if !ok {
		return nil, fmt.Errorf("Failed to append Conjur SSL cert")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}

	client := &http.Client{
		Transport: transport,
	}

	url := fmt.Sprintf("%s/authn/%s/login", applianceURL, url.QueryEscape(account))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create login request. %s", err)
	}
	req.SetBasicAuth(username, string(password))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to send login request. %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to authenticate to conjur. Status code returned '%d'", resp.StatusCode)
	}

	apiKey, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read api key from login request. %s", err)
	}

	return apiKey, nil
}
