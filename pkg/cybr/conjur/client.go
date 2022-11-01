package conjur

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/conjur/authenticators"
	helpersauthn "github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators"
)

var (
	envAccountKey        = "CONJUR_ACCOUNT"
	envApplianceURLKey   = "CONJUR_APPLIANCE_URL"
	envLoginKey          = "CONJUR_AUTHN_LOGIN"
	envAPIKeyKey         = "CONJUR_AUTHN_API_KEY"
	envCertFileKey       = "CONJUR_CERT_FILE"
	envAuthenticatorKey  = "CONJUR_AUTHENTICATOR"
	envAwsTypeKey        = "CONJUR_AWS_TYPE"
	envAuthnServiceIDKey = "CONJUR_AUTHN_SERVICE_ID"
	envSSLVerifyKey      = "CONJUR_SSL_VERIFY"

	envAccount        = os.Getenv(envAccountKey)
	envApplianceURL   = os.Getenv(envApplianceURLKey)
	envLogin          = os.Getenv(envLoginKey)
	envAPIKey         = os.Getenv(envAPIKeyKey)
	envCertFile       = os.Getenv(envCertFileKey)
	envAuthenticator  = os.Getenv(envAuthenticatorKey)
	envAwsType        = os.Getenv(envAwsTypeKey)
	envAuthnServiceID = os.Getenv(envAuthnServiceIDKey)
	envSSLVerify      = os.Getenv(envSSLVerifyKey)
)

func validateEnvironmentConfig(value string, keyName string, errMsg string) string {
	if value == "" {
		errMsg += keyName + ", "
	}
	return errMsg
}

// This method will return an empty conjurapi.Config and error if no conjur environment
// variables are set. If one environment variable is set this will assume that the user
// is attempting to use environment variables and an empty conjurapi.Config is returned with an
// error message of the missing environment variables.
func getClientFromEnvironmentVariable() (*conjurapi.Client, *authn.LoginPair, error) {
	// Do not get client from environment variables because none provided
	if envAccount == "" && envApplianceURL == "" && envLogin == "" && envAPIKey == "" && envAuthenticator == "" {
		return &conjurapi.Client{}, &authn.LoginPair{}, nil
	}

	// Get client from environment variables
	errMsg := ""
	errMsg = validateEnvironmentConfig(envAccount, envAccountKey, errMsg)
	errMsg = validateEnvironmentConfig(envApplianceURL, envApplianceURLKey, errMsg)
	errMsg = validateEnvironmentConfig(envLogin, envLoginKey, errMsg)
	errMsg = validateEnvironmentConfig(envAPIKey, envAPIKeyKey, errMsg)

	// Partial environment variables were provided so return an error
	// with a list of the environment variables that were not provided
	if errMsg != "" {
		return &conjurapi.Client{},
			&authn.LoginPair{},
			fmt.Errorf("environment variable(s) not provided: %s", strings.TrimRight(errMsg, ", "))
	}

	authnURL := helpersauthn.GetAuthURL(envApplianceURL, "authn", "")

	apiKey, err := Login(authnURL, envAccount, envLogin, []byte(envAPIKey), envCertFile)
	if err != nil {
		return &conjurapi.Client{}, &authn.LoginPair{}, err
	}

	config := conjurapi.Config{
		Account:      envAccount,
		ApplianceURL: envApplianceURL,
		SSLCertPath:  envCertFile,
	}
	loginPair := authn.LoginPair{
		Login:  envLogin,
		APIKey: string(apiKey),
	}

	client, err := conjurapi.NewClientFromKey(config, loginPair)
	return client, &loginPair, err
}

func getClientFromAuthenticator() (*conjurapi.Client, *authn.LoginPair, error) {
	errMsg := ""
	errMsg = validateEnvironmentConfig(envAccount, envAccountKey, errMsg)
	errMsg = validateEnvironmentConfig(envApplianceURL, envApplianceURLKey, errMsg)
	errMsg = validateEnvironmentConfig(envLogin, envLoginKey, errMsg)
	errMsg = validateEnvironmentConfig(envAuthenticator, envAuthenticatorKey, errMsg)
	errMsg = validateEnvironmentConfig(envAuthnServiceID, envAuthnServiceIDKey, errMsg)
	if envAuthenticator == "authn-iam" {
		errMsg = validateEnvironmentConfig(envAwsType, envAwsTypeKey, errMsg)
	}

	// Partial environment variables were provided so return an error
	// with a list of the environment variables that were not provided
	if errMsg != "" {
		return &conjurapi.Client{},
			&authn.LoginPair{},
			fmt.Errorf("environment variable(s) not provided: %s", strings.TrimRight(errMsg, ", "))
	}

	envSSLVerifyBool := true
	if envSSLVerify == strings.ToLower("false") || envSSLVerify == strings.ToLower("no") || envSSLVerify == "0" {
		envSSLVerifyBool = false
	}

	config := helpersauthn.Config{
		Account:         envAccount,
		ApplianceURL:    envApplianceURL,
		Login:           envLogin,
		ServiceID:       envAuthnServiceID,
		IgnoreSSLVerify: envSSLVerifyBool,
	}

	authenticator, err := authenticators.GetAuthenticator(envAuthenticator, config)
	if err != nil {
		return &conjurapi.Client{}, &authn.LoginPair{}, err
	}

	client, err := authenticator.Authenticate(config)

	return client, nil, err
}

// GetConjurClient creates a Conjur API client and login pair from environment variables, an authenticator,
// or .netrc & .conjurrc files
func GetConjurClient() (*conjurapi.Client, *authn.LoginPair, error) {
	client, loginPair, err := getClientFromEnvironmentVariable()

	client, loginPair, err = getClientFromAuthenticator()

	// Partial environment variables were provided, assume user is attempting to use environment variables
	if err != nil {
		return &conjurapi.Client{}, &authn.LoginPair{}, err
	}

	// Return client created from environment variables
	if *client != (conjurapi.Client{}) {
		return client, loginPair, nil
	}

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

	loginPair, err = conjurapi.LoginPairFromNetRC(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve credentials from ~/.netrc. %s", err)
	}

	client, err = conjurapi.NewClientFromKey(config, *loginPair)
	return client, loginPair, err
}

// sendConjurAuthenticatedHTTPRequest Send a HTTP request with the conjur session token in the authorization header
// This is used for API endpoints not included in the conjur-api-go SDK (info, whoami, enableauthn)
func sendConjurAuthenticatedHTTPRequest(client *conjurapi.Client, loginPair *authn.LoginPair, url string, method string, body io.Reader) (*http.Response, error) {
	sessionToken, err := client.Authenticate(*loginPair)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate to conjur. %s", err)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request '%s'. %s", url, err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token token=\"%s\"", base64.StdEncoding.EncodeToString(sessionToken)))

	httpClient := client.GetHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request. %s", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 204 {
		return resp, fmt.Errorf("received invalid status code '%d'", resp.StatusCode)
	}

	return resp, nil
}
