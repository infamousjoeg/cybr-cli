package httpjson

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/infamousjoeg/cybr-cli/pkg/logger"
)

type contextKey string

var (
	contextKeyCookies = contextKey("cookies")
)

func (c contextKey) String() string {
	return "httpjson_cookies" + string(c)
}

//Cookies get http cookies from context
func Cookies(ctx context.Context) ([]*http.Cookie, bool) {
	cookies, ok := ctx.Value(contextKeyCookies).(([]*http.Cookie))
	return cookies, ok
}

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

func logRequest(req *http.Request, logger logger.Logger) {
	if logger == nil || !logger.Enabled() {
		return
	}

	logger.Writef("%s %s\n", req.Method, req.URL)

	for key, values := range req.Header {
		for _, value := range values {
			if strings.ToLower(key) == "authorization" {
				logger.Writef("%s: %s\n", key, "*****")
				continue
			}
			logger.Writef("%s: %s\n", key, value)
		}
	}

	if !logger.LogBody() || req.Body == nil {
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	logger.Writeln("")
	body := buf.String()
	logger.Writef("%s\n", body)

	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(body)))
}

func getResponse(ctx context.Context, urls string, method string, token string, body interface{}, insecureTLS bool, logger logger.Logger) (http.Response, error) {
	if insecureTLS {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	}
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(urls)

	cookies, ok := Cookies(ctx)
	if !ok {
		jar.SetCookies(u, nil)
	} else {
		jar.SetCookies(u, cookies)
	}

	httpClient := http.Client{
		Jar:     jar,
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
	req, err := http.NewRequestWithContext(ctx, method, urls, bodyReader)
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

	logRequest(req, logger)
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
func SendRequest(ctx context.Context, url string, method string, token string, body interface{}, insecureTLS bool, logger logger.Logger) (map[string]interface{}, error) {
	res, err := getResponse(ctx, url, method, token, body, insecureTLS, logger)

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
func SendRequestRaw(ctx context.Context, url string, method string, token string, body interface{}, insecureTLS bool, logger logger.Logger) (context.Context, []byte, error) {
	res, err := getResponse(ctx, url, method, token, body, insecureTLS, logger)
	//Read the response body if not nil
	if err != nil && res.Body != nil {
		content, errRead := ioutil.ReadAll(res.Body)
		if errRead != nil {
			return ctx, nil, fmt.Errorf("Failed to read body. %s", errRead)
		}
		newCtx := context.WithValue(ctx, contextKeyCookies, res.Cookies())
		return newCtx, content, err
	}
	return ctx, nil, err
}

// Get a get request and get response as serialized json map[string]interface{}
func Get(url string, token string, insecureTLS bool, logger logger.Logger) (map[string]interface{}, error) {
	ctx := context.TODO()
	response, err := SendRequest(ctx, url, http.MethodGet, token, "", insecureTLS, logger)
	return response, err
}

// Post a post request and get response as serialized json map[string]interface{}
func Post(url string, token string, body interface{}, insecureTLS bool, logger logger.Logger) (map[string]interface{}, error) {
	ctx := context.TODO()
	response, err := SendRequest(ctx, url, http.MethodPost, token, body, insecureTLS, logger)
	return response, err
}

// Put a put request and get response as serialized json map[string]interface{}
func Put(url string, token string, body interface{}, insecureTLS bool, logger logger.Logger) (map[string]interface{}, error) {
	ctx := context.TODO()
	response, err := SendRequest(ctx, url, http.MethodPut, token, body, insecureTLS, logger)
	return response, err
}

// Delete a delete request and get response as serialized json map[string]interface{}
func Delete(url string, token string, insecureTLS bool, logger logger.Logger) (map[string]interface{}, error) {
	ctx := context.TODO()
	response, err := SendRequest(ctx, url, http.MethodDelete, token, "", insecureTLS, logger)
	return response, err
}
