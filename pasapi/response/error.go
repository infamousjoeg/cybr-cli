package response

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/infamousjoeg/pas-api-go/pasapi/logging"
)

// PASError struct contains error code, message, and details
type PASError struct {
	Code    int
	Message string
	Details *PASErrorDetails `json:"error"`
}

// PASErrorDetails struct contains message, code, target, and details mapped to
// an interface
type PASErrorDetails struct {
	Message string
	Code    string
	Target  string
	Details map[string]interface{}
}

// NewPASError creates a PAS error message to return
func NewPASError(resp *http.Response) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cerr := PASError{}
	cerr.Code = resp.StatusCode
	err = json.Unmarshal(body, &cerr)
	if err != nil {
		cerr.Message = strings.TrimSpace(string(body))
	}

	// If the body's empty, use the HTTP status as the message
	if cerr.Message == "" {
		cerr.Message = resp.Status
	}

	return &cerr
}

func (err *PASError) Error() string {
	logging.APILog.Debugf("err.Details: %+v, err.Message: %+v\n", err.Details, err.Message)

	var b strings.Builder

	if err.Message != "" {
		b.WriteString(err.Message + ". ")
	}

	if err.Details != nil && err.Details.Message != "" {
		b.WriteString(err.Details.Message + ".")
	}

	return b.String()
}
