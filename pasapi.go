package pasapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	// BaseURL is a constant exported variable to serve as the default
	// base URL when one is not provided.
	BaseURL = "https://localhost/PasswordVault/api"
)

// Client ...
type Client struct {
	BaseURL      string
	sessionToken string
	HTTPClient   *http.Client
}

// NewClient ...
func NewClient(sessionToken string, baseURL string) *Client {
	return &Client{
		BaseURL:      baseURL,
		sessionToken: sessionToken,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// Content-Type and Body should already be added to req
func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", c.sessionToken)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Try to unmarshall into errorResponse
	if res.StatusCode >= 300 {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	// Unmarshall and populate v
	fullResponse := successResponse{
		Data: v,
	}
	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}

	return nil
}
