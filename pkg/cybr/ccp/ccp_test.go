package ccp_test

import (
	"os"
	"testing"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/ccp"
)

var (
	hostname   = os.Getenv("PAS_HOSTNAME")
	appID      = "cybr-cli-ccp-test"
	safe       = "CLI_ACCOUNTS_TEST"
	object     = "Operating System-UnixSSH-10.0.0.1-test_list"
	clientCert = os.Getenv("CCP_CLIENT_CERT")
	clientKey  = os.Getenv("CCP_CLIENT_PRIVATE_KEY")
)

func TestCCPClientCertSuccess(t *testing.T) {
	query := &ccp.RetrieveAccountQuery{
		AppID:  appID,
		Safe:   safe,
		Object: object,
	}
	request := ccp.RetrieveAccountRequest{
		URL:             hostname,
		IgnoreSSLVerify: false,
		Query:           query,
		ClientCert:      clientCert,
		ClientKey:       clientKey,
	}

	response, err := ccp.RetrieveAccount(request)
	if err != nil {
		t.Errorf("Failed to retrieve account from cyberark using CCP. %s. %v", err, response)
	}
}
