package httpJson

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func getResponse(url string, method string, header http.Header, body string) (http.Response, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	httpClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}
	var bodyReader io.ReadCloser
	var res *http.Response
	bodyReader = ioutil.NopCloser(bytes.NewReader([]byte(body)))

	// create the request
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return *res, fmt.Errorf("Failed to create new request. %s", err)
	}

	// attach the header
	if header == nil {
		header = make(http.Header)
	}

	req.Header = header
	req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Content-Type", "application/json")

	// send request
	res, err = httpClient.Do(req)
	if err != nil {
		return *res, fmt.Errorf("Failed to send request. %s", err)
	}

	if res.StatusCode != 200 {
		return *res, fmt.Errorf("Recieved non-200 status code '%d'", res.StatusCode)
	}

	return *res, err

}

// Send an http request and get response as serialized json map[string]interface{}
func SendRequest(url string, method string, header http.Header, body string) (map[string]interface{}, error) {
	res, err := getResponse(url, method, header, body)
	if err != nil {
		return nil, err
	}

	// Map response body to a map interface
	decoder := json.NewDecoder(res.Body)
	var data map[string]interface{}
	err = decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode response body. %s", err)
	}

	return data, err

}

// Send an http request and get response as byte[]
func SendRequestRaw(url string, method string, header http.Header, body string) ([]byte, error) {
	res, err := getResponse(url, method, header, body)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body. %s", err)
	}
	return content, err

}

// Send an get request and get response as serialized json map[string]interface{}
func Get(url string, header http.Header) (map[string]interface{}, error) {
	response, err := SendRequest(url, http.MethodGet, header, "")
	return response, err
}

// Send an Ppo request and get response as serialized json map[string]interface{}
func Post(url string, header http.Header, body string) (map[string]interface{}, error) {
	response, err := SendRequest(url, http.MethodPost, header, body)
	return response, err
}

// Send an get request and get response as serialized json map[string]interface{}
func Put(url string, header http.Header, body string) (map[string]interface{}, error) {
	response, err := SendRequest(url, http.MethodPut, header, body)
	return response, err
}

// Send an get request and get response as serialized json map[string]interface{}
func Delete(url string, header http.Header) (map[string]interface{}, error) {
	response, err := SendRequest(url, http.MethodDelete, header, "")
	return response, err
}
