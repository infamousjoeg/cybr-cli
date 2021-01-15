package cmd

import (
	"fmt"
	"log"
	"strings"
	"syscall"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	// Username to logon PAS REST API using
	Username string
	// AuthenticationType to be used to logon PAS REST API
	AuthenticationType string
	// InsecureTLS is a boolean value whether to verify TLS or not
	InsecureTLS bool
	// BaseURL to send PAS REST API logon request to
	BaseURL string
)

var logonCmd = &cobra.Command{
	Use:   "logon",
	Short: "Logon to PAS REST API",
	Long: `Authenticate to PAS REST API using the provided authentication type.
	
	Example Usage:
	$ cybr logon -u $USERNAME -a $AUTH_TYPE -b https://pvwa.example.com
	To bypass TLS verification:
	$ cybr logon -u $USERNAME -a $AUTH_TYPE -b https://pvwa.example.com -i`,
	Aliases: []string{"login"},
	Run: func(cmd *cobra.Command, args []string) {
		// Get secret value from STDIN
		fmt.Print("Enter password: ")
		byteSecretVal, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			log.Fatalln("An error occurred trying to read password from " +
				"Stdin. Exiting...")
		}

		client := pasapi.Client{
			BaseURL:     BaseURL,
			AuthType:    AuthenticationType,
			InsecureTLS: InsecureTLS,
		}

		credentials := pasapi.LogonRequest{
			Username: Username,
			Password: string(byteSecretVal),
		}

		err = client.Logon(credentials)
		if err != nil && !strings.Contains(err.Error(), "ITATS542I") {
			log.Fatalf("Failed to Logon to the PVWA. %s", err)
		}

		// if error contains challenge error code, deal with OTPCode here instead and redo client.Logon()
		if err != nil {
			// Get secret value from STDIN
			fmt.Print("Enter one-time passcode: ")
			byteOTPCode, err := terminal.ReadPassword(int(syscall.Stdin))
			credentials.Password = string(byteOTPCode)
			fmt.Println()
			if err != nil {
				log.Fatalln("An error occurred trying to read one-time passcode from " +
					"Stdin. Exiting...")
			}
			err = client.Logon(credentials)
			if err != nil {
				log.Fatalf("Failed to respond to challenge. Possible timeout occurred. %s", err)
			}
		}

		err = client.SetConfig()
		if err != nil {
			log.Fatalf("Failed to create configuration file. %s", err)
			return
		}

		fmt.Printf("Successfully logged onto PAS as user %s.\n", Username)
	},
}

func init() {
	logonCmd.Flags().StringVarP(&Username, "username", "u", "", "Username to logon PAS REST API using")
	logonCmd.MarkFlagRequired("username")
	logonCmd.Flags().StringVarP(&AuthenticationType, "auth-type", "a", "", "Authentication method to logon using [cyberark|ldap|radius-push|radius-append|radius-challenge]")
	logonCmd.MarkFlagRequired("authType")
	logonCmd.Flags().BoolVarP(&InsecureTLS, "insecure-tls", "i", false, "If detected, TLS will not be verified")
	logonCmd.Flags().StringVarP(&BaseURL, "base-url", "b", "", "Base URL to send Logon request to [https://pvwa.example.com]")
	logonCmd.MarkFlagRequired("base-url")
	rootCmd.AddCommand(logonCmd)
}
