package conjur

import (
	"bufio"
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators"
)

var conjurrcTemplate string = `---
account: {{ ACCOUNT }}
plugins: []
appliance_url: {{ APPLIANCE_URL }}
cert_file: "{{ CERT_FILE }}"
authn_url: {{ AUTHN_LDAP_URL }}
`

// GetConjurRcPath returns path to the ~/.conjurrc file os-agnostic
func GetConjurRcPath(homeDir string) string {
	return filepath.FromSlash(fmt.Sprintf("%s/.conjurrc", homeDir))
}

// GetConjurPemPath returns path to the ~/conjur-<account>.pem file os-agnostic
func GetConjurPemPath(homeDir string, account string) string {
	return filepath.FromSlash(fmt.Sprintf("%s/conjur-%s.pem", homeDir, account))
}

func getPem(url string) (string, error) {
	conf := &tls.Config{
		// file deepcode ignore TooPermissiveTrustManager: conjur.pem is not a trusted third-party certificate
		InsecureSkipVerify: true,
	}

	// trim https://
	url = strings.TrimPrefix(url, "https://")
	// If no port is provide default to port 443
	if !strings.Contains(url, ":") {
		url = url + ":443"
	}

	conn, err := tls.Dial("tcp", url, conf)
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve certificate from '%s'. %s", url, err)
	}
	defer conn.Close()

	if len(conn.ConnectionState().PeerCertificates) == 1 {
		return "", fmt.Errorf("Invalid conjur url '%s'. Make sure hostname and port are correct", url)
	}

	pemCert := ""
	for _, cert := range conn.ConnectionState().PeerCertificates {
		pemCert += string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}))
	}

	return pemCert, err
}

func createConjurCert(certFileName string, url string) error {
	// make sure we can get the certificate
	pemCert, err := getPem(url)
	if err != nil {
		return err
	}

	// replace the cert file
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(fmt.Sprintf("Replace certificate file '%s' [y]: ", certFileName))
	text, _ := reader.ReadString('\n')
	answer := strings.TrimSpace(text)
	// overwrite file
	if answer == "" || answer == "y" {
		err = ioutil.WriteFile(certFileName, []byte(pemCert), 0600)
		if err != nil {
			return fmt.Errorf("Failed to write file '%s'. %s", certFileName, err)
		}
	}

	return err
}

func createConjurRcFile(account string, url string, certFileName string, authnLDAP string, conjurrcFileName string) error {
	fmt.Print("Replace ~/.conjurrc file [y]: ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	answer := strings.TrimSpace(text)

	// overwrite ~/.conjurrc file
	if answer == "" || answer == "y" {
		if !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		// create ~/.conjurrrc contents
		conjurrcContent := strings.Replace(conjurrcTemplate, "{{ ACCOUNT }}", account, 1)
		conjurrcContent = strings.Replace(conjurrcContent, "{{ APPLIANCE_URL }}", url, 1)
		conjurrcContent = strings.Replace(conjurrcContent, "{{ CERT_FILE }}", certFileName, 1)
		ldapURL := authenticators.GetAuthURL(url, "authn-ldap", authnLDAP)
		conjurrcContent = strings.Replace(conjurrcContent, "{{ AUTHN_LDAP_URL }}", ldapURL, 1)
		if authnLDAP == "" {
			removeLine := "authn_url: " + ldapURL
			conjurrcContent = strings.Replace(conjurrcContent, removeLine, "", 1)
		}

		err = ioutil.WriteFile(conjurrcFileName, []byte(conjurrcContent), 0600)
		if err != nil {
			return fmt.Errorf("Failed to write file '%s'. %s", conjurrcFileName, err)
		}
	}

	return err
}

func getFieldFromConjurRc(conjurrcFileName string, fieldName string) string {
	file, err := os.Open(conjurrcFileName)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, fieldName+": ") {
			result := strings.SplitN(line, ": ", 2)[1]
			result = strings.Trim(strings.Trim(result, "\n"), "\r")
			return result
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ""
}

// GetHomeDirectory gets the proper user home directory regardless of GOOS
func GetHomeDirectory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Failed to find user's home directory. %s", err)
	}
	return usr.HomeDir, err
}

// GetURLFromConjurRc retrieve conjur url from the ~/.conjurrc file
func GetURLFromConjurRc(conjurrcFileName string) string {
	return getFieldFromConjurRc(conjurrcFileName, "appliance_url")
}

// GetAccountFromConjurRc retrieve conjur account from the ~/.conjurrc file
func GetAccountFromConjurRc(conjurrcFileName string) string {
	return getFieldFromConjurRc(conjurrcFileName, "account")
}

// GetCertFromConjurRc retrieve conjur certificate from the ~/.conjurrc file, is empty if no cert will be used
func GetCertFromConjurRc(conjurrcFileName string) string {
	return strings.Trim(getFieldFromConjurRc(conjurrcFileName, "cert_file"), "\"")
}

// CreateConjurRc creates a ~/.conjurrc file
func CreateConjurRc(account string, url string, selfSignedCert bool, authnLDAP string) error {
	// make sure we can get home directory
	homeDir, err := GetHomeDirectory()
	if err != nil {
		return err
	}

	certFileName := ""
	if selfSignedCert {
		certFileName = GetConjurPemPath(homeDir, account)
		err = createConjurCert(certFileName, url)
		if err != nil {
			return err
		}
	}

	// create the ~/.conjurrc file
	conjurrcFileName := GetConjurRcPath(homeDir)
	err = createConjurRcFile(account, url, certFileName, authnLDAP, conjurrcFileName)

	return err
}
