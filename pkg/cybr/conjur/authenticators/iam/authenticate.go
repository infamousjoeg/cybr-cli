package iam

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func getAuthnURL(authnURL string, account string, login string) string {
	identifier := url.QueryEscape(login)
	return fmt.Sprintf("%s/%s/%s/authenticate", authnURL, account, identifier)
}

func newHTTPSClient(ignoreSSLVerify bool, cert []byte) (*http.Client, error) {
	if ignoreSSLVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		return &http.Client{Transport: tr, Timeout: time.Second * 10}, nil
	}

	// If not certificate provided do not create a certifictae pool
	if len(cert) == 0 {
		return &http.Client{Timeout: time.Second * 10}, nil
	}

	// certificate is provided so create pool and append to TLSClientConfig
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(cert)
	if !ok {
		return nil, fmt.Errorf("Can't append Conjur SSL cert")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}
	return &http.Client{Transport: tr, Timeout: time.Second * 10}, nil
}

// Authenticate to conjur using the authnURL and conjurAuthnRequest
func Authenticate(authnURL string, account string, login string, conjurAuthnRequest string, ignoreSSLVerify bool, cert []byte) ([]byte, error) {
	client, err := newHTTPSClient(ignoreSSLVerify, cert)
	if err != nil {
		return nil, fmt.Errorf("Failed to create a new HTTPS client. %s", err)
	}

	bodyReader := ioutil.NopCloser(bytes.NewReader([]byte(conjurAuthnRequest)))
	url := getAuthnURL(authnURL, account, login)

	response, err := client.Post(url, "", bodyReader)
	if err != nil {
		return nil, fmt.Errorf("Failed to establish connection to Conjur at url '%s'. %s", url, err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to authenticate to Conjur. Received status code '%v'", response.StatusCode)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read Conjur Access Token %s", err)
	}

	return body, err
}
