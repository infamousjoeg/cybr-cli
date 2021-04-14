package conjur

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func getLoginClient(certPath string) (*http.Client, error) {
	// return default client because self-signed certificate will not be used
	client := &http.Client{}
	if certPath == "" {
		return &http.Client{}, nil
	}

	// self-signed certificate needs to be added to the cert pool
	pool := x509.NewCertPool()
	content, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read cert file. %s. %s", certPath, err)
	}

	ok := pool.AppendCertsFromPEM(content)
	if !ok {
		return nil, fmt.Errorf("Failed to append Conjur SSL cert")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}

	client.Transport = transport

	return client, nil
}

// Login to conjur and return an api key in []byte format
func Login(authnURL string, account string, username string, password []byte, certPath string) ([]byte, error) {
	client, err := getLoginClient(certPath)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/login", authnURL, url.QueryEscape(account))
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
