package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/prettyprint"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/util"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/identity"
	identityrequests "github.com/infamousjoeg/cybr-cli/pkg/cybr/identity/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/identity/responses"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/ispss"
	ispssresponses "github.com/infamousjoeg/cybr-cli/pkg/cybr/ispss/responses"
	"github.com/spf13/cobra"
)

// Constants for logon command
const (
	maxAttempts = 10
)

// Global variables for logon command
var (
	Username           string                                 // Username to logon PAS REST API
	AuthenticationType string                                 // Authentication type for PAS REST API
	TenantID           string                                 // Tenant ID for Identity authentication
	InsecureTLS        bool                                   // Boolean to decide whether to verify TLS or not
	BaseURL            string                                 // Base URL to send PAS REST API logon request
	NonInteractive     bool                                   // Flag for non-interactive logon
	Password           string                                 // Password for PAS REST API
	ConcurrentSession  bool                                   // Flag to allow concurrent sessions
	SelectedChallenges []int                                  // Slice of selected challenges for Identity authentication
	AnswerChallenge    identityrequests.AdvanceAuthentication // Answer challenge struct
	StartOobChallenge  identityrequests.AdvanceAuthentication // Start Oob challenge struct
	AnswerOOBChallenge identityrequests.AdvanceAuthentication // Answer Oob challenge struct
	advanceResponse    *responses.Authentication              // Advance authentication response
	platformDiscovery  *ispssresponses.PlatformDiscovery      // Platform discovery response
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

func startAuthIdentity(c pasapi.Client, username string) (*responses.Authentication, error) {
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
	$ cybr logon -u $USERNAME -a identity -b https://example.privilegecloud.cyberark.cloud
	To bypass TLS verification:
	$ cybr logon -u $USERNAME -a $AUTH_TYPE -b https://pvwa.example.com -i`,
	Aliases: []string{"login"},
	Run: func(cmd *cobra.Command, args []string) {
		// Create new client for PAS REST API
		c := pasapi.Client{
			BaseURL:     BaseURL,
			AuthType:    AuthenticationType,
			InsecureTLS: InsecureTLS,
		}

		// Check if auth type is "identity" and get TenantID if true
		if c.AuthType == "identity" {
			platformDiscovery, err := ispss.PlatformDiscovery(c.BaseURL)
			if err != nil {
				log.Fatalf("Failed to get platform discovery. %s", err)
			}
			c.TenantID, err = util.GetSubDomain(platformDiscovery.IdentityUserPortal.API)
			c.BaseURL, err = util.GetSubDomain(platformDiscovery.Pcloud.API)
			if err != nil {
				log.Fatalf("Failed to get tenant ID. %s", err)
			}
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
			startResponse, err := startAuthIdentity(c, Username)
			if err != nil {
				log.Fatalf("%s", err)
			}
			if Verbose {
				prettyprint.PrintColor("cyan", fmt.Sprintf("Start Authentication Response: %+v", startResponse))
			}
			if startResponse.Result.Token != "" {
				c.SessionToken = "Bearer " + startResponse.Result.Token
			}

			// Loop through challenges until c.SessionToken is set
			for attempts := 0; c.SessionToken == "" && attempts < maxAttempts; attempts++ {
			loop:
				// Print challenge number
				prettyprint.PrintColor("yellow", fmt.Sprintf("+ Challenge #%d", attempts+1))

				if startResponse.Result.Challenges[0].Mechanisms[0].PromptSelectMech == "Password" && attempts == 0 {
					// Get password from Stdin
					Password, err = util.ReadPassword()
					if err != nil {
						log.Fatalf("An error occurred trying to read password from Stdin. Exiting")
					}
					// Create AdvanceAuthentication struct
					AnswerChallenge.SessionID = startResponse.Result.SessionID
					AnswerChallenge.MechanismID = startResponse.Result.Challenges[0].Mechanisms[0].MechanismID
					AnswerChallenge.Action = "Answer"
					AnswerChallenge.Answer = Password

					if Verbose {
						prettyprint.PrintColor("cyan", fmt.Sprintf("Session ID: %s", AnswerChallenge.SessionID))
						prettyprint.PrintColor("cyan", fmt.Sprintf("Mechanism ID: %s", AnswerChallenge.MechanismID))
						prettyprint.PrintColor("cyan", fmt.Sprintf("Action: %s", AnswerChallenge.Action))
						prettyprint.PrintColor("cyan", fmt.Sprintf("Answer: %s", AnswerChallenge.Answer))
					}

					// Answer challenge
					advanceResponse, err = identity.AdvanceAuthentication(c, AnswerChallenge)
					if err != nil {
						log.Fatalf("Failed to answer challenge. %s", err)
					}
					if Verbose {
						prettyprint.PrintColor("cyan", fmt.Sprintf("Advance Authentication Response: %+v", advanceResponse))
					}
					if advanceResponse.Result.Token != "" {
						c.SessionToken = "Bearer " + advanceResponse.Result.Token
						break
					}
					if !advanceResponse.Success {
						log.Fatalf("Identity returned unsuccessful response. %s", *advanceResponse.Message)
					}

					attempts++
					goto loop
				}

				if advanceResponse.Result.Summary == "StartNextChallenge" {
					// Ask for user input for challenges
					for i, mechanism := range startResponse.Result.Challenges[1].Mechanisms {
						prettyprint.PrintColor("green", fmt.Sprintf("%d. %s", i+1, mechanism.PromptSelectMech))
					}

					// Get user input
					fmt.Printf("> ")
					var input int
					fmt.Scanln(&input)
					// Keep asking for input until valid
					for input < 1 || input > len(startResponse.Result.Challenges[1].Mechanisms) {
						prettyprint.PrintColor("red", fmt.Sprintf("Please enter a valid number between 1 and %d", len(startResponse.Result.Challenges[1].Mechanisms)))
						fmt.Scanln(&input)
					}
					// Add selected challenge to slice
					SelectedChallenges = append(SelectedChallenges, input)

					// Print selected challenges' PromptSelectMech
					for _, challenge := range SelectedChallenges {
						if Verbose {
							prettyprint.PrintColor("cyan", fmt.Sprintf("Selected: %s", startResponse.Result.Challenges[1].Mechanisms[challenge-1].PromptSelectMech))
						}
						selectedMechanismID := startResponse.Result.Challenges[1].Mechanisms[challenge-1].MechanismID
						selectedAnswerType := startResponse.Result.Challenges[1].Mechanisms[challenge-1].AnswerType

						if strings.HasPrefix(selectedAnswerType, "Start") && strings.HasSuffix(selectedAnswerType, "Oob") {
							StartOobChallenge.SessionID = startResponse.Result.SessionID
							StartOobChallenge.MechanismID = selectedMechanismID
							StartOobChallenge.Action = "StartOOB"

							if Verbose {
								prettyprint.PrintColor("cyan", fmt.Sprintf("Session ID: %s", StartOobChallenge.SessionID))
								prettyprint.PrintColor("cyan", fmt.Sprintf("Mechanism ID: %s", StartOobChallenge.MechanismID))
								prettyprint.PrintColor("cyan", fmt.Sprintf("Action: %s", StartOobChallenge.Action))
							}

							// Answer challenge
							challengeResponse, err := identity.AdvanceAuthentication(c, StartOobChallenge)
							if err != nil {
								log.Fatalf("Failed to answer challenge. %s", err)
							}
							if Verbose {
								prettyprint.PrintColor("cyan", fmt.Sprintf("Advance Authentication Response: %+v", challengeResponse))
							}
							if challengeResponse.Result.Token != "" {
								c.SessionToken = "Bearer " + challengeResponse.Result.Token
								break
							}
							if !challengeResponse.Success {
								log.Fatalf("Identity returned unsuccessful response. %s", *advanceResponse.Message)
							}

							// Get OTP code from Stdin
							fmt.Printf("Enter code: ")
							var code string
							fmt.Scanln(&code)

							AnswerOOBChallenge.SessionID = startResponse.Result.SessionID
							AnswerOOBChallenge.MechanismID = selectedMechanismID
							AnswerOOBChallenge.Action = "Answer"
							AnswerOOBChallenge.Answer = code

							if Verbose {
								prettyprint.PrintColor("cyan", fmt.Sprintf("Session ID: %s", AnswerOOBChallenge.SessionID))
								prettyprint.PrintColor("cyan", fmt.Sprintf("Mechanism ID: %s", AnswerOOBChallenge.MechanismID))
								prettyprint.PrintColor("cyan", fmt.Sprintf("Action: %s", AnswerOOBChallenge.Action))
								prettyprint.PrintColor("cyan", fmt.Sprintf("Answer: %s", AnswerOOBChallenge.Answer))
							}

							// Answer challenge
							answerOOBResponse, err := identity.AdvanceAuthentication(c, AnswerOOBChallenge)
							if err != nil {
								log.Fatalf("Failed to answer challenge. %s", err)
							}
							if Verbose {
								prettyprint.PrintColor("cyan", fmt.Sprintf("Advance Authentication Response: %+v", answerOOBResponse))
							}
							if answerOOBResponse.Result.Token != "" {
								c.SessionToken = "Bearer " + answerOOBResponse.Result.Token
								break
							}
							if advanceResponse.Message != nil {
								log.Fatalf("Identity returned unsuccessful response. %s", *advanceResponse.Message)
							} else {
								log.Fatalf("Identity returned unsuccessful response, but the message is unavailable.")
							}
						}
					}
				}
			}

			// Maximum attempts reached
			if c.SessionToken == "" {
				log.Fatalf("Failed to get non-empty token after %d attempts. Exiting", maxAttempts)
			}

			// Set client config
			err = c.SetConfig()
			if err != nil {
				log.Fatalf("Failed to create configuration file. %s", err)
			}
		}

		// Logon success message
		prettyprint.PrintColor("green", fmt.Sprintf("Successfully logged onto PAS as user %s.", Username))
	},
}

// init function to initialize flags for the 'logon' command
func init() {
	logonCmd.Flags().StringVarP(&Username, "username", "u", "", "Username to logon to PAS REST API")
	logonCmd.MarkFlagRequired("username")
	logonCmd.Flags().StringVarP(&AuthenticationType, "auth-type", "a", "", "Authentication method to logon using [cyberark|ldap|radius]")
	logonCmd.MarkFlagRequired("auth-type")
	logonCmd.Flags().BoolVarP(&InsecureTLS, "insecure-tls", "i", false, "If detected, TLS will not be verified")
	logonCmd.Flags().StringVarP(&BaseURL, "base-url", "b", "", "Base URL to send Logon request to [https://pvwa.example.com]")
	logonCmd.MarkFlagRequired("base-url")
	logonCmd.Flags().BoolVar(&NonInteractive, "non-interactive", false, "If detected, will retrieve the password from the PAS_PASSWORD environment variable")
	logonCmd.Flags().StringVarP(&Password, "password", "p", "", "Password to logon to PAS REST API, only supported when using --non-interactive flag")
	logonCmd.Flags().BoolVar(&ConcurrentSession, "concurrent", false, "If detected, will create a concurrent session to the PAS API")

	// Add 'logon' command to root command
	rootCmd.AddCommand(logonCmd)
}
