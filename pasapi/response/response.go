package response

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/infamousjoeg/pas-api-go/pasapi/logging"
)

func readBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	responseText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseText, err
}

func logResponse(resp *http.Response) {
	req := resp.Request
	logging.APILog.Debugf("%d %s %s %+v", resp.StatusCode, req.Method, req.URL, req.Header)
}

// DataResponse checks the HTTP status of the response. If it's less than
// 300, it returns the response body as a byte array. Otherwise it returns
// a NewPASError.
func DataResponse(resp *http.Response) ([]byte, error) {
	logResponse(resp)
	if resp.StatusCode < 300 {
		return readBody(resp)
	}
	return nil, NewPASError(resp)
}

// SecretDataResponse checks the HTTP status of the response. If it's less than
// 300, it returns the response body as a stream. Otherwise it returns
// a NewPASError.
func SecretDataResponse(resp *http.Response) (io.ReadCloser, error) {
	logResponse(resp)
	if resp.StatusCode < 300 {
		return resp.Body, nil
	}
	return nil, NewPASError(resp)
}

// JSONResponse checks the HTTP status of the response. If it's less than
// 300, it returns the response body as JSON. Otherwise it returns
// a NewPASError.
func JSONResponse(resp *http.Response, obj interface{}) error {
	logResponse(resp)
	if resp.StatusCode < 300 {
		body, err := readBody(resp)
		if err != nil {
			return err
		}
		return json.Unmarshal(body, obj)
	}
	return NewPASError(resp)
}

// EmptyResponse checks the HTTP status of the response. If it's less than
// 300, it returns without an error. Otherwise it returns
// a NewPASError.
func EmptyResponse(resp *http.Response) error {
	logResponse(resp)
	if resp.StatusCode < 300 {
		return nil
	}
	return NewPASError(resp)
}
