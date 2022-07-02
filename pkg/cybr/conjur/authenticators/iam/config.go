package iam

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// ConjurAwsType is the environment variable that holds the aws type
	ConjurAwsType = "CONJUR_AWS_TYPE"
	// ConjurAccount is the environment variable that holds the account name
	ConjurAccount = "CONJUR_ACCOUNT"
	// ConjurApplianceURL is the environment variable that holds the appliance URL
	ConjurApplianceURL = "CONJUR_APPLIANCE_URL"
	// ConjurAuthnURL is the environment variable that holds the authn URL
	ConjurAuthnURL = "CONJUR_AUTHN_URL"
	// ConjurAuthnLogin is the environment variable that holds the authn login
	ConjurAuthnLogin = "CONJUR_AUTHN_LOGIN"
	// ConjurAccessTokenPath is the environment variable that holds the access token path
	ConjurAccessTokenPath = "CONJUR_ACCESS_TOKEN_PATH"
	// ConjurIgnoreSSLVerify is the environment variable that holds the ignore ssl verify status
	ConjurIgnoreSSLVerify = "CONJUR_IGNORE_SSL_VERIFY"
	// ConjurRefresh is the environment variable that holds the refresh status
	ConjurRefresh = "CONJUR_REFRESH"
	// ConjurSSLCertificate is the environment variable that holds the ssl certificate
	ConjurSSLCertificate = "CONJUR_SSL_CERTIFICATE"
	// ConjurCertFile is the environment variable that holds the certificate file path
	ConjurCertFile = "CONJUR_CERT_FILE"

	// FlagAwsType is the flag that holds the aws type
	FlagAwsType = "aws-name"
	// FlagAccount is the flag that holds the account name
	FlagAccount = "account"
	// FlagApplianceURL is the flag that holds the appliance URL
	FlagApplianceURL = "URL"
	// FlagLogin is the flag that holds the login
	FlagLogin = "login"
	// FlagAuthnURL is the flag that holds the authn URL
	FlagAuthnURL = "authn-URL"
	// FlagTokenPath is the flag that holds the access token path
	FlagTokenPath = "token-path"
	// FlagSecretID is the flag that holds the secret id
	FlagSecretID = "secret"
	// FlagSilence is the flag that holds the silence status
	FlagSilence = "silence"
	// FlagIgnoreSSLVerify is the flag that holds the ignore ssl verify status
	FlagIgnoreSSLVerify = "ignore-ssl-verify"
	// FlagRefresh is the flag that holds the refresh status
	FlagRefresh = "refresh"
	// FlagCertFile is the flag that holds the certificate file path
	FlagCertFile = "cert-file"

	// DescriptionAwsType is the description of the aws type
	DescriptionAwsType = "AWS Resource type name. Environment variable equivalent '" + ConjurAwsType + "'. e.g. ec2, lambda, ecs"
	// DescriptionAccount is the description of the account name
	DescriptionAccount = "The Conjur account. Environment variable equivalent '" + ConjurAccount + "'. e.g. company, etc"
	// DescriptionApplianceURL is the description of the appliance URL
	DescriptionApplianceURL = "The URL to the Conjur instance. Environment variable equivalent '" + ConjurApplianceURL + "'. e.g. https://conjur.com"
	// DescriptionLogin is the description of the login
	DescriptionLogin = "Conjur login that will be used. Environment variable equivalent '" + ConjurAuthnLogin + "'. e.g. host/6634674884744/iam-role-name"
	// DescriptionAuthnURL is the description of the authn URL
	DescriptionAuthnURL = "URL Conjur will be authenticating to. Environment variable equivalent '" + ConjurAuthnURL + "'. e.g. https://conjur.com/authn-iam/global"
	// DescriptionTokenPath is the description of the access token path
	DescriptionTokenPath = "Write the access token to this file. Environment variable equivalent '" + ConjurAccessTokenPath + "'. e.g. /path/to/access-token.json"
	// DescriptionSecretID is the description of the secret id
	DescriptionSecretID = "Retrieve a specific secret from Conjur. e.g. db/postgres/username"
	// DescriptionSilence is the description of the silence status
	DescriptionSilence = "Silence debug and info messages"
	// DescriptionIgnoreSSLVerify is the description of the ignore ssl verify status
	DescriptionIgnoreSSLVerify = "WARNING: Do not verify the SSL certificate provided by Conjur server. THIS SHOULD ONLY BE USED FOR POC"
	// DescriptionRefresh is the description of the refresh status
	DescriptionRefresh = "Continously run and retrieve the Conjur access token every 6 min"
	// DescriptionCertFile is the description of the certificate file path
	DescriptionCertFile = "The Conjur certificate chain file. Environment variable equivalent '" + ConjurCertFile + "'. e.g. /etc/conjur.pem"
)

// Config holds the configuration for the Conjur authenticator
type Config struct {
	AWSName         string
	Account         string
	ApplianceURL    string
	Login           string
	AuthnURL        string
	IgnoreSSLVerify bool
	Certificate     string
	CertificatePath string

	// If AccessTokenPath & SecretID is not provided then print access token to stdout
	// If only AccessTokenPath is provided then write access token to file
	// If only SecretID is provided then print secret value to stdout
	// If AccessTokenPath & SecretID is provided then write access token to file and print secret value to stdout
	// FetchSeed will retrieve the seed from the master seed service and will write it to a file
	AccessTokenPath string
	SecretID        string
	Silence         bool

	// For retrying on failure, default is 5
	Retry     int
	RetryWait int64
	Refresh   bool
}

func (c Config) getCertificate() ([]byte, error) {
	if c.Certificate != "" {
		return []byte(c.Certificate), nil
	}

	if c.CertificatePath != "" {
		return ioutil.ReadFile(c.CertificatePath)
	}
	return nil, nil
}

func newHTTPSClient(ignoreSSLVerify bool, cert []byte) (*http.Client, error) {
	// If not certificate provided do not create a certifictae pool
	if cert == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: ignoreSSLVerify},
		}
		return &http.Client{Transport: tr, Timeout: time.Second * 10}, nil
	}

	// certificate is provided so create pool and append to TLSClientConfig
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(cert)
	if !ok {
		return nil, fmt.Errorf("Can't append Conjur SSL cert")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}
	return &http.Client{Transport: tr, Timeout: time.Second * 10}, nil
}

// GetConfig will default to using environment variables if flag is not provided.
// If environment variable and flag is provided then the flag will override the environment variable
func GetConfig() (Config, error) {
	// mandatory properties
	awsName := flag.String(FlagAwsType, os.Getenv(ConjurAwsType), DescriptionAwsType)
	account := flag.String(FlagAccount, os.Getenv(ConjurAccount), DescriptionAccount)
	applianceURL := flag.String(FlagApplianceURL, os.Getenv(ConjurApplianceURL), DescriptionApplianceURL)
	login := flag.String(FlagLogin, os.Getenv(ConjurAuthnLogin), DescriptionLogin)
	authnURL := flag.String(FlagAuthnURL, os.Getenv(ConjurAuthnURL), DescriptionAuthnURL)

	// optional properties
	tokenPath := flag.String(FlagTokenPath, os.Getenv(ConjurAccessTokenPath), DescriptionTokenPath)
	secretID := flag.String(FlagSecretID, "", DescriptionSecretID)
	silence := flag.Bool(FlagSilence, false, DescriptionSilence)
	certFile := flag.String(FlagCertFile, os.Getenv(ConjurCertFile), DescriptionCertFile)

	ignoreStr := strings.ToLower(os.Getenv(ConjurIgnoreSSLVerify))
	ignoreDefault := false
	if ignoreStr == "yes" || ignoreStr == "true" {
		ignoreDefault = true
	}
	ignoreSSLVerify := flag.Bool(FlagIgnoreSSLVerify, ignoreDefault, DescriptionIgnoreSSLVerify)

	refreshStr := strings.ToLower(os.Getenv(ConjurRefresh))
	refreshDefault := false
	if refreshStr == "yes" || refreshStr == "true" {
		refreshDefault = true
	}
	refresh := flag.Bool(FlagRefresh, refreshDefault, DescriptionRefresh)

	flag.Parse()

	// Validate mandatory config properties
	if *awsName == "" {
		return Config{}, fmt.Errorf("%s, %s", ConjurAwsType, FlagAwsType)
	}

	if *account == "" {
		return Config{}, fmt.Errorf("%s, %s", ConjurAccount, FlagAccount)
	}

	if *applianceURL == "" {
		return Config{}, fmt.Errorf("%s, %s", ConjurApplianceURL, FlagApplianceURL)
	}

	if *login == "" {
		return Config{}, fmt.Errorf("%s, %s", ConjurAuthnLogin, FlagLogin)
	}

	if *authnURL == "" {
		return Config{}, fmt.Errorf("%s, %s", ConjurAuthnURL, FlagAuthnURL)
	}

	config := Config{
		AWSName:         *awsName,
		Account:         *account,
		ApplianceURL:    *applianceURL,
		Login:           *login,
		AuthnURL:        *authnURL,
		AccessTokenPath: *tokenPath,
		SecretID:        *secretID,
		Silence:         *silence,
		IgnoreSSLVerify: *ignoreSSLVerify,
		Retry:           5,
		RetryWait:       60,
		Refresh:         *refresh,
		CertificatePath: *certFile,
	}
	return config, nil
}
