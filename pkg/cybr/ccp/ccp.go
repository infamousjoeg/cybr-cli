package ccp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	httpJson "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
)

// RetrieveAccountRequest request used to retrieve account from the CCP
type RetrieveAccountRequest struct {
	URL             string
	IgnoreSSLVerify bool
	ClientCert      string
	ClientKey       string
	Query           *RetrieveAccountQuery
}

// RetrieveAccountQuery represents valid query parameters when listing accounts
type RetrieveAccountQuery struct {
	AppID                       string `query_key:"AppID"`
	Safe                        string `query_key:"Safe"`
	Folder                      string `query_key:"Folder"`
	Object                      string `query_key:"Object"`
	UserName                    string `query_key:"UserName"`
	Address                     string `query_key:"Address"`
	Database                    string `query_key:"Database"`
	PolicyID                    string `query_key:"PolicyID"`
	ConnectionTimeout           string `query_key:"ConnectionTimeout"`
	Query                       string `query_key:"Query"`
	QueryFormat                 string `query_key:"QueryFormat"`
	FailRequestOnPasswordChange bool   `query_key:"FailRequestOnPasswordChange"`
}

// RetrieveAccount retrieves an account for CCP and returns results via map
func RetrieveAccount(request RetrieveAccountRequest) (map[string]string, error) {
	query := httpJson.GetURLQuery(request.Query)
	url := ccpURL(request.URL, query)
	cert, useClientCert, err := clientCertificate(request.ClientCert, request.ClientKey)
	if err != nil {
		return map[string]string{}, err
	}

	return sendHTTPRequest(url, request.IgnoreSSLVerify, useClientCert, cert)
}

func ccpURL(url string, query string) string {
	return fmt.Sprintf("%s/AIMWebService/api/Accounts%s", url, query)
}

func clientCertificate(clientCert string, clientKey string) (tls.Certificate, bool, error) {
	// If either the cert of key are empty then return an empty tls certificate
	if clientCert == "" || clientKey == "" {
		return tls.Certificate{}, false, nil
	}

	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return tls.Certificate{}, false, fmt.Errorf("Failed to load client cert and key from files '%s' and '%s'. %s", clientCert, clientKey, err)
	}

	return cert, true, nil
}

func sendHTTPRequest(url string, insecureSkipVerify bool, useClientCert bool, cert tls.Certificate) (map[string]string, error) {
	tlsConfig := &tls.Config{
		Renegotiation:      tls.RenegotiateOnceAsClient,
		InsecureSkipVerify: insecureSkipVerify,
	}

	if useClientCert {
		tlsConfig.Certificates = []tls.Certificate{cert}
		tlsConfig.BuildNameToCertificate()
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to send request to url '%s'. %s", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		// if we fail to read body when recieveing an invalid status code, ignore.
		errorMessage, _ := streamToStringMap(resp.Body)
		jsonErrorMessage, _ := json.Marshal(errorMessage)
		return errorMessage, fmt.Errorf("Invalid response from CCP url '%s'. Status Code: %s. %s", url, resp.Status, jsonErrorMessage)
	}

	return streamToStringMap(resp.Body)
}

func streamToStringMap(stream io.Reader) (map[string]string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return map[string]string{}, err
	}
	body := buf.Bytes()

	jsonMap := make(map[string]string)
	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		return map[string]string{}, fmt.Errorf("Failed to unmarshal returned body into string map. %s", err)
	}

	return jsonMap, nil
}
