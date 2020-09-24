package cmd

import (
	"fmt"
	"log"
	"syscall"

	pasapi "github.com/infamousjoeg/pas-api-go/pkg/cybr/api"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	// Username to logon PAS REST API using
	Username string
	// AuthenticationType to be used to logon PAS REST API
	AuthenticationType string
	// BaseURL to send PAS REST API logon request to
	BaseURL string
)

var logonCmd = &cobra.Command{
	Use:   "logon",
	Short: "Logon to PAS REST API",
	Long: `Authenticate to PAS REST API using the provided authentication type.
	
	Example Usage:
	$ cybr logon -u $USERNAME -a $AUTH_TYPE -b https://pvwa.example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get secret value from STDIN
		fmt.Println("Enter password: ")
		byteSecretVal, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalln("An error occurred trying to read password from " +
				"Stdin. Exiting...")
		}

		client := pasapi.Client{
			BaseURL:  BaseURL,
			AuthType: AuthenticationType,
		}

		credentials := pasapi.LogonRequest{
			Username: Username,
			Password: string(byteSecretVal),
		}

		err = client.Logon(credentials)
		if err != nil {
			log.Fatalf("Failed to Logon to the PVWA. %s", err)
			return
		}

		fmt.Printf("Successfully logged onto PAS as user %s.\n", Username)
	},
}

func init() {
	logonCmd.Flags().StringVarP(&Username, "username", "u", "", "Username to logon PAS REST API using")
	logonCmd.Flags().StringVarP(&AuthenticationType, "authType", "a", "", "Authentication method to logon using")
	logonCmd.Flags().StringVarP(&BaseURL, "base-url", "b", "", "Base URL to send Logon request to")
	rootCmd.AddCommand(logonCmd)
}
