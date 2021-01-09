package httpjson

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func bodyToBytes(body interface{}) ([]byte, error) {
	if body == nil {
		return []byte(""), nil
	}

	content, err := json.Marshal(body)
	if err != nil {
		return []byte(""), err
	}
	return content, nil
}

func getResponse(url string, method string, token string, body interface{}, insecureTLS bool) (http.Response, error) {
	if insecureTLS {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	}
	httpClient := http.Client{
		Timeout: time.Second * 30, // Maximum of 30 secs
	}
	var bodyReader io.ReadCloser
	var res *http.Response

	content, err := bodyToBytes(body)
	if err != nil {
		return *res, err
	}

	bodyReader = ioutil.NopCloser(bytes.NewReader(content))

	// create the request
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return *res, fmt.Errorf("Failed to create new request. %s", err)
	}

	// attach the header
	req.Header = make(http.Header)
	req.Header.Add("Content-Type", "application/json")
	// if token is provided, add header Authorization
	if token != "" {
		req.Header.Add("Authorization", token)
	}

	// send request
	res, err = httpClient.Do(req)
	if err != nil {
		return http.Response{}, fmt.Errorf("Failed to send request. %s", err)
	}

	if res.StatusCode >= 300 {
		return *res, fmt.Errorf("Received non-200 status code '%d'", res.StatusCode)
	}

	return *res, err
}

// SendRequest is an http request and get response as serialized json map[string]interface{}
func SendRequest(url string, method string, token string, body interface{}, insecureTLS bool) (map[string]interface{}, error) {
	res, err := getResponse(url, method, token, body, insecureTLS)

	if err != nil && strings.Contains(err.Error(), "Failed to send request") {
		return nil, err
	}
	if res.StatusCode == 204 {
		return nil, nil
	}

	// Map response body to a map interface
	decoder := json.NewDecoder(res.Body)
	var data map[string]interface{}
	decodeError := decoder.Decode(&data)

	// No error and no body returned
	if decodeError == io.EOF {
		return nil, nil
	}

	if decodeError != nil {
		return nil, fmt.Errorf("Failed to decode response body. %s", err)
	}

	return data, err
}

// SendRequestRaw is an http request and get response as byte[]
func SendRequestRaw(url string, method string, token string, body interface{}, insecureTLS bool) ([]byte, error) {
	res, err := getResponse(url, method, token, body, insecureTLS)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body. %s", err)
	}
	return content, err
}

// Get a get request and get response as serialized json map[string]interface{}
func Get(url string, token string, insecureTLS bool) (map[string]interface{}, error) {
	response, err := SendRequest(url, http.MethodGet, token, "", insecureTLS)
	return response, err
}

// Post a post request and get response as serialized json map[string]interface{}
func Post(url string, token string, body interface{}, insecureTLS bool) (map[string]interface{}, error) {
	response, err := SendRequest(url, http.MethodPost, token, body, insecureTLS)
	return response, err
}

// Put a put request and get response as serialized json map[string]interface{}
func Put(url string, token string, body interface{}, insecureTLS bool) (map[string]interface{}, error) {
	response, err := SendRequest(url, http.MethodPut, token, body, insecureTLS)
	return response, err
}

// Delete a delete request and get response as serialized json map[string]interface{}
func Delete(url string, token string, insecureTLS bool) (map[string]interface{}, error) {
	response, err := SendRequest(url, http.MethodDelete, token, "", insecureTLS)
	return response, err
}
