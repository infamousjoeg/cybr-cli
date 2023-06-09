package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/util"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/identity"
	identityrequests "github.com/infamousjoeg/cybr-cli/pkg/cybr/identity/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/identity/responses"
	"github.com/spf13/cobra"
)

// Global variables for logon command
var (
	Username           string // Username to logon PAS REST API
	AuthenticationType string // Authentication type for PAS REST API
	TenantID           string // Tenant ID for Identity authentication
	InsecureTLS        bool   // Boolean to decide whether to verify TLS or not
	BaseURL            string // Base URL to send PAS REST API logon request
	NonInteractive     bool   // Flag for non-interactive logon
	Password           string // Password for PAS REST API
	ConcurrentSession  bool   // Flag to allow concurrent sessions
	SelectedChallenges []int  // Slice of selected challenges for Identity authentication
)

func logonToPAS(c pasapi.Client, username, password string, nonInteractive, concurrentSession bool) error {
	var err error
	// Check if non-interactive flag is not provided and password is not empty
	if !nonInteractive && password != "" {
		return fmt.Errorf("An error occured because --non-interactive must be provided when using --password flag")
	}
	// If the execution is not non-interactive, ask the user to input password
	if !nonInteractive {
		password, err = util.ReadPassword()
		if err != nil {
			return fmt.Errorf("An error occurred trying to read password from Stdin. Exiting")
		}
	}
	// Check if password is empty
	if password == "" {
		return fmt.Errorf("Provided password is empty")
	}
	// Create credentials for logon
	credentials := requests.Logon{
		Username:          username,
		Password:          password,
		ConcurrentSession: concurrentSession,
	}
	// Logon to the PAS REST API
	err = c.Logon(credentials)
	if err != nil && !strings.Contains(err.Error(), "ITATS542I") {
		return fmt.Errorf("Failed to Logon to the PVWA. %s", err)
	}
	// Deal with OTPCode here if error contains challenge error code and redo client.Logon()
	if err != nil {
		// Get OTP code from Stdin
		credentials, err = util.ReadOTPcode(credentials)
		err = c.Logon(credentials)
		if err != nil {
			return fmt.Errorf("Failed to respond to challenge. Possible timeout occurred. %s", err)
		}
	}
	// Set client config
	err = c.SetConfig()
	if err != nil {
		return fmt.Errorf("Failed to create configuration file. %s", err)
	}
	return nil
}

func startAuthIdentity(c pasapi.Client, username string) (*responses.StartAuthentication, error) {
	// Create credentials for logon
	startAuth := identityrequests.StartAuthentication{
		User:     username,
		TenantID: c.TenantID,
		Version:  "1.0",
	}

	// Start authentication
	response, err := identity.StartAuthentication(c, startAuth)
	if err != nil {
		return nil, fmt.Errorf("Failed to start authentication. %s", err)
	}
	if response.Success != true {
		return nil, fmt.Errorf("Identity returned unsuccessful response. %s", *response.Message)
	}

	return response, nil
}

// logonCmd represents the 'logon' command for PAS REST API
var logonCmd = &cobra.Command{
	Use:   "logon",
	Short: "Logon to PAS REST API",
	Long: `Authenticate to PAS REST API using the provided authentication type.
	
	Example Usage:
	$ cybr logon -u $USERNAME -a $AUTH_TYPE -b https://pvwa.example.com
	Logon to Privilege Cloud REST API:
	$ cybr logon -u $USERNAME -a identity -t xxx1234 -b https://example.privilegecloud.cyberark.cloud
	To bypass TLS verification:
	$ cybr logon -u $USERNAME -a $AUTH_TYPE -b https://pvwa.example.com -i`,
	Aliases: []string{"login"},
	Run: func(cmd *cobra.Command, args []string) {
		// Create new client for PAS REST API
		c := pasapi.Client{
			BaseURL:     BaseURL,
			AuthType:    AuthenticationType,
			TenantID:    TenantID,
			InsecureTLS: InsecureTLS,
		}
		// Check if auth type is "identity" and no tenant id is provided
		if c.AuthType == "identity" && c.TenantID == "" {
			log.Fatalf("An error occured because --tenant-id must be provided when using --auth-type identity")
		}

		// Get password from environment variable PAS_PASSWORD
		Password := os.Getenv("PAS_PASSWORD")

		// Handle authentication depending on auth type
		if c.AuthType != "identity" {
			err := logonToPAS(c, Username, Password, NonInteractive, ConcurrentSession)
			if err != nil {
				log.Fatalf("%s", err)
			}
			// Handle Identity authentication
		} else {
			// Start authentication
			response, err := startAuthIdentity(c, Username)
			if err != nil {
				log.Fatalf("%s", err)
			}

			// If token is present, set config and exit
			if response.Result.Token != "" {
				// Set token
				c.SessionToken = response.Result.Token
				// Set client config
				err = c.SetConfig()
				if err != nil {
					log.Fatalf("Failed to create configuration file. %s", err)
				}
				// Logon success message
				fmt.Printf("Successfully logged onto PAS as user %s.\n", Username)
				return
			}

			if response.Result.Challenges[0].Mechanisms[0].PromptSelectMech == "Password" {
				// Get password from Stdin
				Password, err = util.ReadPassword()
				if err != nil {
					log.Fatalf("An error occurred trying to read password from Stdin. Exiting")
				}
				AnswerChallenge1 := identityrequests.AdvanceAuthentication{
					SessionID:   response.Result.SessionID,
					MechanismID: response.Result.Challenges[0].Mechanisms[0].MechanismID,
					Action:      "Answer",
					Answer:      Password,
				}
				fmt.Printf("Session ID: %s\n", AnswerChallenge1.SessionID)
				fmt.Printf("Mechanism ID: %s\n", AnswerChallenge1.MechanismID)
				fmt.Printf("Action: %s\n", AnswerChallenge1.Action)
				fmt.Printf("Answer: %s\n", AnswerChallenge1.Answer)
			} else {
				// Ask for user input for challenges
				for _, challenge := range response.Result.Challenges {
					for i, mechanism := range challenge.Mechanisms {
						fmt.Printf("%d. %s\n", i+1, mechanism.PromptSelectMech)
					}

					// Get user input
					fmt.Printf("> ")
					var input int
					fmt.Scanln(&input)
					// Keep asking for input until valid
					for input < 1 || input > len(response.Result.Challenges[0].Mechanisms) {
						fmt.Printf("Please enter a valid number between 1 and %d\n", len(response.Result.Challenges[0].Mechanisms))
						fmt.Scanln(&input)
					}
					// Add selected challenge to slice
					SelectedChallenges = append(SelectedChallenges, input)
				}

				// Print selected challenges' PromptSelectMech
				for _, challenge := range SelectedChallenges {
					fmt.Printf("Selected: %s\n", response.Result.Challenges[0].Mechanisms[challenge-1].PromptSelectMech)
				}
			}
		}

		// Logon success message
		fmt.Printf("Successfully logged onto PAS as user %s.\n", Username)
	},
}

// init function to initialize flags for the 'logon' command
func init() {
	logonCmd.Flags().StringVarP(&Username, "username", "u", "", "Username to logon to PAS REST API")
	logonCmd.MarkFlagRequired("username")
	logonCmd.Flags().StringVarP(&AuthenticationType, "auth-type", "a", "", "Authentication method to logon using [cyberark|ldap|radius]")
	logonCmd.MarkFlagRequired("auth-type")
	logonCmd.Flags().StringVarP(&TenantID, "tenant-id", "t", "", "The ID of the Identity tenant to authenticate to [e.g. xxx1234]")
	logonCmd.Flags().BoolVarP(&InsecureTLS, "insecure-tls", "i", false, "If detected, TLS will not be verified")
	logonCmd.Flags().StringVarP(&BaseURL, "base-url", "b", "", "Base URL to send Logon request to [https://pvwa.example.com]")
	logonCmd.MarkFlagRequired("base-url")
	logonCmd.Flags().BoolVar(&NonInteractive, "non-interactive", false, "If detected, will retrieve the password from the PAS_PASSWORD environment variable")
	logonCmd.Flags().StringVarP(&Password, "password", "p", "", "Password to logon to PAS REST API, only supported when using --non-interactive flag")
	logonCmd.Flags().BoolVar(&ConcurrentSession, "concurrent", false, "If detected, will create a concurrent session to the PAS API")

	// Add 'logon' command to root command
	rootCmd.AddCommand(logonCmd)
}
