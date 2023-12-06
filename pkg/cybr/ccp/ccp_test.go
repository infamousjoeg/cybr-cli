package ccp_test

import (
	"os"
	"testing"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/ccp"
)

var (
	hostname   = os.Getenv("CCP_HOSTNAME")
	appID      = "cybr-cli-ccp-test"
	safe       = "PIN-APP-CYBRCLI-TEST"
	object     = "Operating System-PL-WIN-DOMAIN-ADMIN-10.0.0.1-test-new"
	clientCert = os.Getenv("CCP_CLIENT_CERT")
	clientKey  = os.Getenv("CCP_CLIENT_PRIVATE_KEY")
)

func writeCertsToFile(clientCertContent string, clientKeyContent string) (string, string, error) {
	certFilePath := os.TempDir() + "/client.crt"
	keyFilePath := os.TempDir() + "/client.key"

	err := os.WriteFile(certFilePath, []byte(clientCertContent), 0644)
	if err != nil {
		return "", "", err
	}
	err = os.WriteFile(keyFilePath, []byte(clientKeyContent), 0644)
	if err != nil {
		return "", "", err
	}

	return certFilePath, keyFilePath, nil

}

func TestCCPClientCertSuccess(t *testing.T) {
	clientCertPath, clientKeyPath, err := writeCertsToFile(clientCert, clientKey)
	if err != nil {
		t.Fatalf("Failed to write the client cert and key. %s", err)
	}

	query := &ccp.RetrieveAccountQuery{
		AppID:  appID,
		Safe:   safe,
		Object: object,
	}
	request := ccp.RetrieveAccountRequest{
		URL:             hostname,
		IgnoreSSLVerify: false,
		Query:           query,
		ClientCert:      clientCertPath,
		ClientKey:       clientKeyPath,
	}

	response, err := ccp.RetrieveAccount(request)
	if err != nil {
		t.Errorf("Failed to retrieve account from cyberark using CCP. %s. %v", err, response)
	}

	// Clean up the files
	os.Remove(clientCertPath)
	os.Remove(clientKeyPath)
}
