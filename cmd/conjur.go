package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/conjur"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/prettyprint"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

const stdinErrMsg = "Failed to read from stdin."

var (
	// Account conjur account
	Account string

	// AuthnLDAP Authenticator Service ID
	AuthnLDAP string

	// PolicyBranch branch policy is being loaded into
	PolicyBranch string

	// PolicyFilePath path to policy file
	PolicyFilePath string

	// VariableID variable ID of a secret
	VariableID string

	// NoNewLine no new line when printing secret
	NoNewLine bool

	// SecretValue variable secret value
	SecretValue string

	// ServiceID used for enabling authenticator
	ServiceID string

	// Kind resource kind variable, policy, user, host, group, etc
	Kind string

	// InspectResources inspect the resources and provide more info per resource
	InspectResources bool
)

func isInputFromPipe() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func loadPolicyFile(policyBranch string, policyFilePath string, policyMode conjurapi.PolicyMode) {
	if policyFilePath == "" {
		log.Fatal("Policy file path is required")
	}

	client, _, err := conjur.GetConjurClient()
	if err != nil {
		log.Fatalf("Failed to initialize conjur client. %s", err)
	}

	file, err := os.Open(policyFilePath)
	if err != nil {
		log.Fatalf("Failed to read policy file '%s'. %s", policyFilePath, err)
	}

	response, err := client.LoadPolicy(policyMode, policyBranch, bufio.NewReader(file))
	if err != nil {
		log.Fatalf("Failed to load policy. %v. %s", response, err)
	}
	prettyprint.PrintJSON(response)
}

func loadPolicyPipe(policyBranch, policyContent string, policyMode conjurapi.PolicyMode) {
	client, _, err := conjur.GetConjurClient()
	if err != nil {
		log.Fatalf("Failed to initialize conjur client. %s", err)
	}

	response, err := client.LoadPolicy(policyMode, policyBranch, strings.NewReader(policyContent))
	if err != nil {
		log.Fatalf("Failed to load policy. %v. %s", response, err)
	}
	prettyprint.PrintJSON(response)
}

func removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatalf("Failed to remove file '%s'. %s", path, err)
	}
}

func readPassword() []byte {
	// Convert Password variable to byte array
	byteSecretVal := []byte(Password)

	// If password is not provided, prompt for password
	if len(byteSecretVal) == 0 {
		fmt.Print("Enter password: ")
		byteSecretVal, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("An error occurred trying to read password from " +
				"Stdin. Exiting...")
		}
		fmt.Println()
		return byteSecretVal
	}

	return byteSecretVal
}

var conjurCmd = &cobra.Command{
	Use:   "conjur",
	Short: "Conjur actions",
	Long:  `Perform actions on conjur`,
}

var conjurLogonCmd = &cobra.Command{
	Use:   "logon",
	Short: "Logon to Conjur",
	Long: `Authenticate to Conjur using API Key or password
	
	Example Usage:
	$ cybr conjur logon -a account -b https://conjur.example.com -l admin
	$ cybr conjur logon -a account -b https://conjur.example.com -l serviceAccountUser --authn-ldap`,
	Aliases: []string{"login"},
	Run: func(cmd *cobra.Command, args []string) {
		byteSecretVal := readPassword()

		homeDir, err := conjur.GetHomeDirectory()
		if err != nil {
			log.Fatalf("%s\n", err)
		}

		netrcPath := conjur.GetNetRcPath(homeDir)

		// certPath remains empty if not using self-signed-cert
		certPath := ""
		if InsecureTLS {
			certPath = conjur.GetConjurPemPath(homeDir, Account)
		}

		err = conjur.CreateConjurRc(Account, BaseURL, InsecureTLS, AuthnLDAP)
		if err != nil {
			log.Fatalf("Failed to create ~/.conjurrc file. %s\n", err)
		}

		authnURL := authenticators.GetAuthURL(BaseURL, "authn", "")
		if AuthnLDAP != "" {
			authnURL = authenticators.GetAuthURL(BaseURL, "authn-ldap", AuthnLDAP)
		}

		apiKey, err := conjur.Login(authnURL, Account, Username, byteSecretVal, certPath)
		if err != nil {
			log.Fatalf("Failed to login and retrieve api key. %s", err)
		}

		err = conjur.CreateNetRc(Username, string(apiKey))
		if err != nil {
			log.Fatalf("Failed to create ~/.netrc file. %s\n", err)
		}

		config := conjurapi.Config{
			Account:      Account,
			ApplianceURL: BaseURL,
			NetRCPath:    netrcPath,
			SSLCertPath:  certPath,
		}

		loginPair := authn.LoginPair{
			Login:  Username,
			APIKey: string(apiKey),
		}

		client, err := conjurapi.NewClientFromKey(config, loginPair)
		_, err = client.Authenticate(loginPair)
		if err != nil {
			log.Fatalf("Failed to authenticate to conjur. %s", err)
		}

		fmt.Println("Successfully logged into conjur")
	},
}

var conjurNonInteractiveLogonCmd = &cobra.Command{
	Use:   "logon-non-interactive",
	Short: "Logon to Conjur non-interactively",
	Long: `Authenticate to Conjur using API Key or password non-interactively 
	
	Example Usage:
	$ cybr conjur logon-non-interactive`,
	Aliases: []string{"login-non-interactive"},
	Run: func(cmd *cobra.Command, args []string) {
		client, loginPair, err := conjur.GetConjurClient()
		if err != nil {
			log.Fatalf("Failed to initialize conjur client. %s", err)
		}

		_, err = client.Authenticate(*loginPair)
		if err != nil {
			log.Fatalf("Failed to authenticate to conjur. %s", err)
		}

		fmt.Println("Successfully logged into conjur")

	},
}

var conjurLogoffCmd = &cobra.Command{
	Use:   "logoff",
	Short: "Logoff to Conjur",
	Long: `Logoff to conjur and remove the ~/.netrc and ~/.conjurrc files
	
	Example Usage:
	$ cybr conjur logoff`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := conjur.GetHomeDirectory()
		if err != nil {
			log.Fatalf("%s\n", err)
		}

		netrcPath := fmt.Sprintf("%s/.netrc", homeDir)
		conjurrcPath := fmt.Sprintf("%s/.conjurrc", homeDir)

		removeFile(netrcPath)
		removeFile(conjurrcPath)

		fmt.Println("Logged off conjur")
	},
}

var conjurAppendPolicyCmd = &cobra.Command{
	Use:   "append-policy",
	Short: "Append policy to conjur",
	Long: `Adds data to the existing Conjur policy. Deletions are not allowed. 
	Any policy objects that exist on the server but are omitted from the policy file will not be deleted and any explicit deletions in the policy file will result in an error.  
	
	Example Usage:
	$ cybr conjur append-policy --branch root --file ./path/to/root.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		if isInputFromPipe() {
			// Read from stdin
			policy, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("%s %s", stdinErrMsg, err)
			}
			loadPolicyPipe(PolicyBranch, string(policy), conjurapi.PolicyModePost)
		} else {
			loadPolicyFile(PolicyBranch, PolicyFilePath, conjurapi.PolicyModePost)
		}
	},
}

var conjurUpdatePolicyCmd = &cobra.Command{
	Use:   "update-policy",
	Short: "Update policy to conjur",
	Long: `Modifies an existing Conjur policy. Data may be explicitly deleted using the !delete, !revoke, and !deny statements. 
	Unlike “replace” mode, no data is ever implicitly deleted.
	
	Example Usage:
	$ cybr conjur update-policy --branch root --file ./path/to/root.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		if isInputFromPipe() {
			// Read from stdin
			policy, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("%s %s", stdinErrMsg, err)
			}
			loadPolicyPipe(PolicyBranch, string(policy), conjurapi.PolicyModePut)
		} else {
			loadPolicyFile(PolicyBranch, PolicyFilePath, conjurapi.PolicyModePatch)
		}
	},
}

var conjurReplacePolicyCmd = &cobra.Command{
	Use:   "replace-policy",
	Short: "Replace policy to conjur",
	Long: `Loads or replaces a Conjur policy document.
	Any policy data which already exists on the server but is not explicitly specified in the new policy file will be deleted.
	
	Example Usage:
	$ cybr conjur replace-policy --branch root --file ./path/to/root.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		if isInputFromPipe() {
			// Read from stdin
			policy, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("%s %s", stdinErrMsg, err)
			}
			loadPolicyPipe(PolicyBranch, string(policy), conjurapi.PolicyModePut)
		} else {
			loadPolicyFile(PolicyBranch, PolicyFilePath, conjurapi.PolicyModePut)
		}
	},
}

var conjurGetSecretCmd = &cobra.Command{
	Use:   "get-secret",
	Short: "Retrieve secret from conjur",
	Long: `Fetches the value of a secret from the specified Variable. 
	The latest version will be retrieved unless the version parameter is specified. 
	The twenty most recent secret versions are retained.
	
	Example Usage:
	$ cybr conjur get-secret -i id/to/variable`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _, err := conjur.GetConjurClient()
		if err != nil {
			log.Fatalf("Failed to initialize conjur client. %s", err)
		}

		content, err := client.RetrieveSecret(VariableID)
		if err != nil {
			log.Fatalf("Failed to retrieve secret variable '%s'. %s", VariableID, err)
		}

		padding := "\n"
		if NoNewLine {
			padding = ""
		}
		fmt.Printf("%s%s", string(content), padding)
	},
}

var conjurSetSecretCmd = &cobra.Command{
	Use:   "set-secret",
	Short: "Set secret in conjur",
	Long: `Sets a secret value for the specified Variable.
	
	Example Usage:
	$ cybr conjur set-secret -i id/to/variable -v "P@$$word"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _, err := conjur.GetConjurClient()
		if err != nil {
			log.Fatalf("Failed to initialize conjur client. %s", err)
		}

		err = client.AddSecret(VariableID, SecretValue)
		if err != nil {
			log.Fatalf("Failed to set secret variable '%s'. %s", VariableID, err)
		}
	},
}

var conjurEnableAuthnCmd = &cobra.Command{
	Use:   "enable-authn",
	Short: "Enable a conjur authenticator",
	Long: `Enables a conjur authenticator.
	
	Example Usage:
	$ cybr conjur enable-authn -s authn-iam/prod`,
	Run: func(cmd *cobra.Command, args []string) {
		err := conjur.EnableAuthenticator(ServiceID)
		if err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Printf("Successfully enabled authenticator '%s'\n", ServiceID)
	},
}

var conjurInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get info about conjur",
	Long: `Get info about conjur.
	
	Example Usage:
	$ cybr conjur info`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := conjur.Info()
		if err != nil {
			log.Fatalf("%s", err)
		}
		prettyprint.PrintJSON(result)
	},
}

var conjurWhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Get current user info logged into Conjur",
	Long: `Get current user information logged into Conjur.
	
	Example Usage:
	$ cybr conjur whoami`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := conjur.Whoami()
		if err != nil {
			log.Fatalf("%s", err)
		}
		prettyprint.PrintJSON(result)
	},
}

var conjurListResourcesCmd = &cobra.Command{
	Use:   "list",
	Short: "List conjur resources",
	Long: `Lists resources within an organization account.
	
	Example Usage:
	$ cybr conjur list --kind variable --search prod`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _, err := conjur.GetConjurClient()
		if err != nil {
			log.Fatalf("Failed to initialize conjur client. %s", err)
		}

		fmt.Printf("client: %+v\n", client)
		fmt.Printf("error: %+v\n", err)

		filter := conjurapi.ResourceFilter{
			Kind:   Kind,
			Search: Search,
			Limit:  Limit,
			Offset: Offset,
		}

		resources, err := client.Resources(&filter)
		if err != nil {
			log.Fatalf("Failed to list resources. %s", err)
		}

		if InspectResources {
			prettyprint.PrintJSON(resources)
			return
		}

		// Just display resource ids
		ids := []string{}
		for _, r := range resources {
			ids = append(ids, r["id"].(string))
		}
		prettyprint.PrintJSON(ids)
	},
}

var conjurRotateAPIKeyCmd = &cobra.Command{
	Use:   "rotate-api-key",
	Short: "Rotate my or other host/user api key",
	Long: `Replaces the API key of another role you can update with a new, securely random API key. The new API key is returned as the response body.


	
	Example Usage:
	$ cybr conjur rotate-api-key
	$ cybr conjur rotate-api-key -l admin
	$ cybr conjur rotate-api-key -l host/some/application`,
	Run: func(cmd *cobra.Command, args []string) {
		client, loginPair, err := conjur.GetConjurClient()
		if err != nil {
			log.Fatalf("Failed to initialize conjur client. %s", err)
		}

		login := Username
		if login == "" {
			login = loginPair.Login
		}

		// Create fully qualified name, if not starts with `host/` then assume user
		if strings.HasPrefix(login, "host/") {
			login = strings.Replace(login, "host/", "host:", 1)
		} else {
			login = "user:" + login
		}

		login = client.GetConfig().Account + ":" + login

		newAPIKey, err := client.RotateAPIKey(login)
		if err != nil {
			log.Fatalf("Failed to rotate api key for '%s'. %s", login, err)
		}

		fmt.Println(string(newAPIKey))
	},
}

func init() {
	// Logon command
	conjurLogonCmd.Flags().StringVarP(&Username, "login", "l", "", "Conjur login name")
	conjurLogonCmd.MarkFlagRequired("login")
	conjurLogonCmd.Flags().StringVarP(&Password, "password", "p", "", "Conjur password")
	conjurLogonCmd.Flags().StringVarP(&Account, "account", "a", "", "Conjur account")
	conjurLogonCmd.MarkFlagRequired("account")
	conjurLogonCmd.Flags().StringVarP(&BaseURL, "base-url", "b", "", "Conjur appliance URL")
	conjurLogonCmd.MarkFlagRequired("base-url")
	conjurLogonCmd.Flags().StringVarP(&AuthnLDAP, "authn-ldap", "", "", "Uses provided Service ID to configure LDAP Authentication")
	conjurLogonCmd.Flags().BoolVar(&InsecureTLS, "self-signed", false, "Retrieve and use self-signed certificate when sending requests to the Conjur API")

	// append-policy
	conjurAppendPolicyCmd.Flags().StringVarP(&PolicyBranch, "branch", "b", "", "The policy branch in which policy is being loaded")
	conjurAppendPolicyCmd.MarkFlagRequired("branch")
	conjurAppendPolicyCmd.Flags().StringVarP(&PolicyFilePath, "file", "f", "", "The policy file that will be loaded into the branch")

	// update-policy
	conjurUpdatePolicyCmd.Flags().StringVarP(&PolicyBranch, "branch", "b", "", "The policy branch in which policy is being loaded")
	conjurUpdatePolicyCmd.MarkFlagRequired("branch")
	conjurUpdatePolicyCmd.Flags().StringVarP(&PolicyFilePath, "file", "f", "", "The policy file that will be loaded into the branch")
	conjurUpdatePolicyCmd.MarkFlagRequired("file")

	// replace-policy
	conjurReplacePolicyCmd.Flags().StringVarP(&PolicyBranch, "branch", "b", "", "The policy branch in which policy is being loaded")
	conjurReplacePolicyCmd.MarkFlagRequired("branch")
	conjurReplacePolicyCmd.Flags().StringVarP(&PolicyFilePath, "file", "f", "", "The policy file that will be loaded into the branch")
	conjurReplacePolicyCmd.MarkFlagRequired("file")

	// retrieve-secret
	conjurGetSecretCmd.Flags().StringVarP(&VariableID, "id", "i", "", "The variable ID containing the secret")
	conjurGetSecretCmd.MarkFlagRequired("ID")
	conjurGetSecretCmd.Flags().BoolVarP(&NoNewLine, "no-new-line", "n", false, "Remove new line")

	// set-secret
	conjurSetSecretCmd.Flags().StringVarP(&VariableID, "id", "i", "", "The variable ID being updated")
	conjurSetSecretCmd.MarkFlagRequired("ID")
	conjurSetSecretCmd.Flags().StringVarP(&SecretValue, "secret-value", "v", "", "The new value of the secret")
	conjurSetSecretCmd.MarkFlagRequired("secret-value")

	// enable-authn
	conjurEnableAuthnCmd.Flags().StringVarP(&ServiceID, "service-id", "s", "", "The authenticator service ID. e.g. authn-iam/prod or authn-k8s/k8s-cluster-1")
	conjurEnableAuthnCmd.MarkFlagRequired("service-id")

	// list
	conjurListResourcesCmd.Flags().StringVarP(&Kind, "kind", "k", "", "Narrows results to only resources of that kind")
	conjurListResourcesCmd.Flags().StringVarP(&Search, "search", "s", "", "Narrows results to those pertaining to the search query")
	conjurListResourcesCmd.Flags().IntVarP(&Limit, "limit", "l", 0, "Maximum number of returned resources")
	conjurListResourcesCmd.Flags().IntVarP(&Offset, "offset", "o", 0, "Index to start returning results from for pagination")
	conjurListResourcesCmd.Flags().BoolVarP(&InspectResources, "inspect", "i", false, "Show full object information")

	// rotate-api-key
	conjurRotateAPIKeyCmd.Flags().StringVarP(&Username, "login", "l", "", "Replaces the API key of another role you can update with a new, securely random API key. The new API key is returned as the response body. e.g. admin, host/someApp")

	conjurCmd.AddCommand(conjurLogonCmd)
	conjurCmd.AddCommand(conjurNonInteractiveLogonCmd)
	conjurCmd.AddCommand(conjurAppendPolicyCmd)
	conjurCmd.AddCommand(conjurUpdatePolicyCmd)
	conjurCmd.AddCommand(conjurReplacePolicyCmd)
	conjurCmd.AddCommand(conjurGetSecretCmd)
	conjurCmd.AddCommand(conjurSetSecretCmd)
	conjurCmd.AddCommand(conjurEnableAuthnCmd)
	conjurCmd.AddCommand(conjurInfoCmd)
	conjurCmd.AddCommand(conjurWhoamiCmd)
	conjurCmd.AddCommand(conjurListResourcesCmd)
	conjurCmd.AddCommand(conjurRotateAPIKeyCmd)
	conjurCmd.AddCommand(conjurLogoffCmd)
	rootCmd.AddCommand(conjurCmd)
}
