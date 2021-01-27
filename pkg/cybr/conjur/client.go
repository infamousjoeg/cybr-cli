package conjur

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
)

// GetNetRcPath returns path to the ~/.netrc file os-agnostic
func GetNetRcPath(homeDir string) string {
	return filepath.FromSlash(fmt.Sprintf("%s/.netrc", homeDir))
}

// GetConjurClient create conjur client and login pair for ~/.conjurrc and ~/.netrc
func GetConjurClient() (*conjurapi.Client, *authn.LoginPair, error) {
	homeDir, err := GetHomeDirectory()
	if err != nil {
		return nil, nil, fmt.Errorf("%s", err)
	}

	netrcPath := GetNetRcPath(homeDir)
	conjurrcPath := GetConjurRcPath(homeDir)

	account := GetAccountFromConjurRc(conjurrcPath)
	baseURL := GetURLFromConjurRc(conjurrcPath)
	certPath := GetCertFromConjurRc(conjurrcPath)

	config := conjurapi.Config{
		Account:      account,
		ApplianceURL: baseURL,
		SSLCertPath:  certPath,
		NetRCPath:    netrcPath,
	}

	loginPair, err := conjurapi.LoginPairFromNetRC(config)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to retrieve credentials from ~/.netrc. %s", err)
	}

	client, err := conjurapi.NewClientFromKey(config, *loginPair)
	return client, loginPair, err
}

// sendConjurAuthenticatedHTTPRequest Send a HTTP request with the conjur session token in the authorization header
func sendConjurAuthenticatedHTTPRequest(client *conjurapi.Client, loginPair *authn.LoginPair, url string, method string, body io.Reader) (*http.Response, error) {
	sessionToken, err := client.Authenticate(*loginPair)
	if err != nil {
		return nil, fmt.Errorf("Failed to authenticate to conjur. %s", err)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request '%s'. %s", url, err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token token=\"%s\"", base64.StdEncoding.EncodeToString(sessionToken)))

	httpClient := client.GetHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to send HTTP request. %s", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 204 {
		return resp, fmt.Errorf("Recieved invalid status code '%d'", resp.StatusCode)
	}

	return resp, nil
}
