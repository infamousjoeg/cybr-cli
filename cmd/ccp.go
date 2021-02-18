package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/ccp"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/prettyprint"
	"github.com/spf13/cobra"
)

var (
	// IgnoreSSLVerify Ignore SSL Verification
	IgnoreSSLVerify bool
	// ClientCert path to the client cert file
	ClientCert string
	// ClientKey path to the client private key file
	ClientKey string
	// Folder in which account resides
	Folder string
	// ObjectName in which account resides
	ObjectName string
	// Database in which account resides
	Database string
	// ConnectionTimeout to wait for CCP
	ConnectionTimeout string
	// Query for the account
	Query string
	// QueryFormat query format being used
	QueryFormat string
	// FailRequestOnPasswordChange if password is currently in a change process
	FailRequestOnPasswordChange bool
	// Field that will be parsed and returned from the account
	Field string
)

var ccpCmd = &cobra.Command{
	Use:   "ccp",
	Short: "CCP actions",
	Long: `All actions that can be performed with the Central Credential Provider.
	
	Example Usage:
	Get an account: $ cybr ccp get-account -b https://ccp.company.local -i AppID -s SafeName -o ObjectName -f Username`,
}

var ccpGetAccountCmd = &cobra.Command{
	Use:   "get-account",
	Short: "Get account from CCP",
	Long: `Get account from the CCP.


	Example Usage:
	$ cybr ccp get-account -b https://ccp.company.local -i AppID -s SafeName -o ObjectName -f Username`,
	Run: func(cmd *cobra.Command, args []string) {
		if ClientCert != "" && ClientKey == "" {
			log.Fatalf("Client certificate was provided with no client private key")
		}
		if ClientKey != "" && ClientCert == "" {
			log.Fatalf("Client private key was provided with no client certificate")
		}

		query := &ccp.RetrieveAccountQuery{
			AppID:                       AppID,
			Safe:                        SafeName,
			Folder:                      Folder,
			Object:                      ObjectName,
			UserName:                    Username,
			Address:                     Address,
			Database:                    Database,
			PolicyID:                    PlatformID,
			ConnectionTimeout:           ConnectionTimeout,
			Query:                       Query,
			QueryFormat:                 QueryFormat,
			FailRequestOnPasswordChange: FailRequestOnPasswordChange,
		}

		request := ccp.RetrieveAccountRequest{
			URL:             BaseURL,
			IgnoreSSLVerify: IgnoreSSLVerify,
			ClientCert:      ClientCert,
			ClientKey:       ClientKey,
			Query:           query,
		}

		account, err := ccp.RetrieveAccount(request)
		if err != nil {
			log.Fatalf("%s", err)
		}

		if Field == "" {
			prettyprint.PrintJSON(account)
			return
		}

		for key, value := range account {
			if strings.ToLower(key) == strings.ToLower(Field) {
				fmt.Println(value)
				return
			}
		}

		log.Fatalf("Failed to parse field '%s' from account returned", Field)
	},
}

// RetrieveAccountQuery represents valid query parameters when listing accounts
type RetrieveAccountQuery struct {
	AppID                       string `query_key:"AppID"`
	Safe                        string `query_key:"Safe"`
	Folder                      string `query_key:"Folder"`
	Object                      string `query_key:"Object"`
	UserName                    string `query_key:"UserName"`
	Address                     string `query_key:"Address"`
	Database                    string `query_key:"Database"`
	PolicyID                    string `query_key:"PolicyID"`
	ConnectionTimeout           string `query_key:"ConnectionTimeout"`
	Query                       string `query_key:"Query"`
	QueryFormat                 string `query_key:"QueryFormat"`
	FailRequestOnPasswordChange bool   `query_key:"FailRequestOnPasswordChange"`
}

func init() {
	// Logon command
	ccpGetAccountCmd.Flags().StringVarP(&BaseURL, "base-url", "b", "", "CCP Base url. e.g. https://ccp.company.local")
	ccpGetAccountCmd.MarkFlagRequired("base-url")
	ccpGetAccountCmd.Flags().StringVarP(&AppID, "app-id", "i", "", "CCP application ID")
	ccpGetAccountCmd.MarkFlagRequired("app-id")
	ccpGetAccountCmd.Flags().BoolVar(&IgnoreSSLVerify, "ignore-ssl-verification", false, "Ignore SSL verification when connecting to CCP server")
	ccpGetAccountCmd.Flags().StringVarP(&ClientCert, "client-cert", "c", "", "Path to the client certificate file")
	ccpGetAccountCmd.Flags().StringVarP(&ClientKey, "client-key", "k", "", "Path to the client private key file")
	ccpGetAccountCmd.Flags().StringVarP(&SafeName, "safe", "s", "", "The safe in which the account resides")
	ccpGetAccountCmd.Flags().StringVar(&Folder, "folder", "", "The folder in which the account resides")
	ccpGetAccountCmd.Flags().StringVarP(&ObjectName, "object-name", "o", "", "The account object name")
	ccpGetAccountCmd.Flags().StringVarP(&Username, "username", "u", "", "The account's username")
	ccpGetAccountCmd.Flags().StringVarP(&Address, "address", "a", "", "The account's address")
	ccpGetAccountCmd.Flags().StringVarP(&Database, "database", "d", "", "The account's database")
	ccpGetAccountCmd.Flags().StringVarP(&PlatformID, "platform-id", "p", "", "The account's platform ID")
	ccpGetAccountCmd.Flags().StringVarP(&ConnectionTimeout, "connection-timeout", "t", "", "Timeout period for CCP to retrieve the account")
	ccpGetAccountCmd.Flags().StringVarP(&Query, "query", "q", "", "The query to perform when retrieveing an account")
	ccpGetAccountCmd.Flags().StringVar(&QueryFormat, "query-format", "", "The query format. Possible values are: Exact or Regexp")
	ccpGetAccountCmd.Flags().BoolVar(&FailRequestOnPasswordChange, "fail-request-on-password-change", false, "Fail the request if password is currenlty being changed")
	ccpGetAccountCmd.Flags().StringVarP(&Field, "field", "f", "", "To return only one field so JSON parsing is not required")

	ccpCmd.AddCommand(ccpGetAccountCmd)
	rootCmd.AddCommand(ccpCmd)
}
