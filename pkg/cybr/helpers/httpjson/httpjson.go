package httpjson

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
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

// Cookies returns the cookies from the context
func Cookies(ctx context.Context) ([]*http.Cookie, bool) {
	cookies, ok := ctx.Value(contextKeyCookies).([]*http.Cookie)
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

	req.Body = io.NopCloser(bytes.NewReader([]byte(body)))
}

func getResponse(ctx context.Context, identity bool, urls string, method string, token string, body interface{}, insecureTLS bool, logger logger.Logger) (http.Response, error) {
	if insecureTLS {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	}

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(urls)
	cookies, ok := Cookies(ctx)
	if ok {
		jar.SetCookies(u, cookies)
	} else {
		jar.SetCookies(u, nil)
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

	bodyReader = io.NopCloser(bytes.NewReader(content))

	// create the request
	req, err := http.NewRequestWithContext(ctx, method, urls, bodyReader)
	if err != nil {
		return *res, fmt.Errorf("Failed to create new request. %s", err)
	}

	// attach the header
	req.Header = make(http.Header)
	req.Header.Add("Content-Type", "application/json")
	if identity {
		req.Header.Add("X-IDAP-NATIVE-CLIENT", "true")
	}
	// if token is provided, add header Authorization
	if token != "" {
		if identity {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		} else {
			req.Header.Add("Authorization", token)
		}
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
func SendRequest(ctx context.Context, identity bool, url string, method string, token string, body interface{}, insecureTLS bool, logger logger.Logger) (map[string]interface{}, error) {
	res, err := getResponse(ctx, identity, url, method, token, body, insecureTLS, logger)

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
func SendRequestRaw(ctx context.Context, identity bool, url string, method string, token string, body interface{}, insecureTLS bool, logger logger.Logger) (context.Context, []byte, error) {
	res, err := getResponse(ctx, identity, url, method, token, body, insecureTLS, logger)
	if err != nil {
		return ctx, nil, err
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return ctx, nil, fmt.Errorf("Failed to read body. %v", err)
	}

	// Check if the response status code is 500 and look for the error message ITATS542I
	if res.StatusCode == 500 && bytes.Contains(content, []byte("ITATS542I")) {
		newCtx := context.WithValue(ctx, contextKeyCookies, res.Cookies())
		return newCtx, content, err
	} else if res.StatusCode >= 300 {
		return ctx, nil, fmt.Errorf("Received non-200 status code '%d'", res.StatusCode)
	}

	newCtx := context.WithValue(ctx, contextKeyCookies, res.Cookies())
	return newCtx, content, err

	// res, err := getResponse(ctx, identity, url, method, token, body, insecureTLS, logger)
	// fmt.Printf("SRR Response: %s\n", res.Body)
	// fmt.Printf("SRR Response Error:	%s\n", err)
	// if err != nil && res.Body != nil {
	// 	content, errRead := io.ReadAll(res.Body)
	// 	if errRead != nil {
	// 		return ctx, nil, fmt.Errorf("Failed to read body. %s", errRead)
	// 	}
	// 	newCtx := context.WithValue(ctx, contextKeyCookies, res.Cookies())
	// 	return newCtx, content, err
	// }
	// return ctx, nil, err
}

// SendRequestRawWithHeaders is an http request and get response as byte[]
func SendRequestRawWithHeaders(url, method string, headers http.Header, body interface{}, insecureTLS bool, logger logger.Logger) ([]byte, error) {
	if insecureTLS {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	}
	httpClient := http.Client{
		Timeout: time.Second * 30, // Maximum of 30 secs
	}

	var res *http.Response

	content, err := bodyToBytes(body)
	if err != nil {
		return []byte(""), err
	}

	// create the request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(content))
	if err != nil {
		return []byte(""), fmt.Errorf("Failed to create new request. %s", err)
	}

	// attach the header
	req.Header = headers

	logRequest(req, logger)

	// send request
	res, err = httpClient.Do(req)
	if err != nil {
		return []byte(""), fmt.Errorf("Failed to send request. %s", err)
	}

	if res.StatusCode >= 300 {
		return []byte(""), fmt.Errorf("Received non-200 status code '%d'", res.StatusCode)
	}

	content, err = io.ReadAll(res.Body)
	if err != nil {
		return []byte(""), fmt.Errorf("Failed to read body. %s", err)
	}
	return content, err
}

// Get a get request and get response as serialized json map[string]interface{}
func Get(identity bool, url string, token string, insecureTLS bool, logger logger.Logger) (map[string]interface{}, error) {
	ctx := context.TODO()

	response, err := SendRequest(ctx, identity, url, http.MethodGet, token, "", insecureTLS, logger)
	return response, err
}

// Post a post request and get response as serialized json map[string]interface{}
func Post(identity bool, url string, token string, body interface{}, insecureTLS bool, logger logger.Logger) (map[string]interface{}, error) {
	ctx := context.TODO()

	response, err := SendRequest(ctx, identity, url, http.MethodPost, token, body, insecureTLS, logger)
	return response, err
}

// Put a put request and get response as serialized json map[string]interface{}
func Put(identity bool, url string, token string, body interface{}, insecureTLS bool, logger logger.Logger) (map[string]interface{}, error) {
	ctx := context.TODO()

	response, err := SendRequest(ctx, identity, url, http.MethodPut, token, body, insecureTLS, logger)
	return response, err
}

// Delete a delete request and get response as serialized json map[string]interface{}
func Delete(identity bool, url string, token string, insecureTLS bool, logger logger.Logger) (map[string]interface{}, error) {
	ctx := context.TODO()

	response, err := SendRequest(ctx, identity, url, http.MethodDelete, token, "", insecureTLS, logger)
	return response, err
}
