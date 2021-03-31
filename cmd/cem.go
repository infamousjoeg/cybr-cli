package cmd

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/quincycheng/cem-api-go/cem"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	// CemOrganization Organization tenant name
	CemOrganization string

	// CemNonInteractive NonInteractive logon
	CemNonInteractive bool

	// CemPlatform Platform Name
	CemPlatform string
	// CemAccountID Account ID
	CemAccountID string
	// CemEntityID Entity ID
	CemEntityID string

	// CemNonFullAdmin non-full Admin only
	CemNonFullAdmin bool
	// CemNonShadowAdmin non-sadow Admin only
	CemNonShadowAdmin bool
	// CemFullAdmin Full Admin only
	CemFullAdmin bool
	// CemShadowAdmin Shadow Admin
	CemShadowAdmin bool

	// CemNextToken Next Token
	CemNextToken string

	// CemSessionTokenPath path to session token file
	CemSessionTokenPath string = "/.cybr/cem.config"

	// CemEnvAPIKey environment variable of CEM API Key for non-interfactive logon
	CemEnvAPIKey string = "CEM_APIKEY"
)

var cemCmd = &cobra.Command{
	Use:   "cem",
	Short: "CEM actions",
	Long: `All actions that can be performed with the Cloud Entitlements Manager.
	
	Example Usage:
	Get remediations on an entity: 
	$ cybr cem get-remediations -p PLATFORM -a ACCOUNT_ID -e ENTITY_ID`,
}

var cemLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to CEM REST API",
	Long: `Authenticate to Cloud Entitlements Manager REST API.
	
	Example Usage:
	$ cybr cem login -a $ORGANIZATION
	
	For non-interactive logon:
	$ export ` + CemEnvAPIKey + `=<Your CEM access key>
	$ cybr cem login -a $ORGANIZATION --non-interactive`,
	Run: func(cmd *cobra.Command, args []string) {
		apikey := os.Getenv(CemEnvAPIKey)

		if !CemNonInteractive {
			// Get secret value from STDIN
			fmt.Print("Enter API Key: ")
			byteSecretVal, err := terminal.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			if err != nil {
				log.Fatalln("An error occurred trying to read API key from " +
					"Stdin. Exiting...")
			}
			apikey = string(byteSecretVal)
		}

		if apikey == "" {
			log.Fatalf("Provided API Key is empty")
		}

		token, err := cem.Login(CemOrganization, apikey)
		if err != nil {
			log.Fatal(err)
		}

		err = SaveToken(token)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Successfully logged onto CEM (organization: %s)\n", CemOrganization)
	},
}

// getUserHomeDir Get the Home directory of the current user
func getUserHomeDir() (string, error) {
	// Get user home directory
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not read user home directory for OS. %s", err)
	}
	return userHome, nil
}

// SaveToken saving token as file on the local filesystem
func SaveToken(token string) error {
	// Get user home directory
	userHome, err := getUserHomeDir()
	if err != nil {
		return fmt.Errorf("ACL error. %s", err)
	}

	// Check if .cybr directory already exists, create if not
	if _, err = os.Stat(userHome + "/.cybr"); os.IsNotExist(err) {
		// Create .cybr folder in user home directory
		err = os.Mkdir(userHome+"/.cybr", 0766)
		if err != nil {
			return fmt.Errorf("could not create folder %s/.cybr on local file system. %s", userHome, err)
		}
	}

	// Check for config file and remove if existing
	if _, err = os.Stat(userHome + CemSessionTokenPath); !os.IsNotExist(err) {
		err = os.Remove(userHome + CemSessionTokenPath)
		if err != nil {
			return fmt.Errorf("could not remove existing %s%s file. %s", userHome, CemSessionTokenPath, err)
		}
	}
	// Create config file in user home directory
	dataFile, err := os.Create(userHome + CemSessionTokenPath)
	if err != nil {
		return fmt.Errorf("could not create configuration file at %s%s. %s", userHome, CemSessionTokenPath, err)
	}

	// serialize the data
	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(token)

	dataFile.Close()

	return nil
}

// GetToken file from local filesystem and read
func GetToken() (string, error) {

	// Get user home directory
	userHome, err := getUserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ACL error. %s", err)
	}

	// open data file
	dataFile, err := os.Open(userHome + CemSessionTokenPath)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve token file at %s. %s", CemSessionTokenPath, err)
	}

	dataDecoder := gob.NewDecoder(dataFile)
	result := ""
	err = dataDecoder.Decode(&result)
	if err != nil {
		return result, fmt.Errorf("failed to decode token file at .cybr/config. %s", err)
	}

	dataFile.Close()

	return result, nil
}

var cemGetAccountsCmd = &cobra.Command{
	Use:   "get-accounts",
	Short: "Get Accounts",
	Long: `Retrieve workspaces grouped by platforms

	Example Usage:
	$ cybr cem get-accounts`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := GetToken()
		if err != nil {
			log.Fatalf("Failed to retrieve token file at %s. %s\n", CemSessionTokenPath, err)
		}

		cemResult, err := cem.GetAccounts(token)
		if err != nil {
			log.Fatalln(err)
			return
		}

		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, []byte(cemResult), "", "  ")
		if error != nil {
			log.Fatal("JSON parse error: ", error)
		}
		fmt.Println(prettyJSON.String())

	},
}

var cemGetRemediationsCmd = &cobra.Command{
	Use:   "get-remediations",
	Short: "Get Entity Remediations",
	Long: `Retrieve Remediations for an entity

	Example Usage:
	$ cybr cem get-remediations -p PLATFORM -a ACCOUNT_ID -e ENTITY_ID`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := GetToken()
		if err != nil {
			log.Fatalf("Failed to retrieve token file at %s. %s\n", CemSessionTokenPath, err)
		}

		cemQuery := &cem.EntityQuery{
			Platform:  CemPlatform,
			AccountId: CemAccountID,
			EntityId:  CemEntityID,
		}

		cemResult, err := cem.GetEntityRecommendations(token, cemQuery)
		if err != nil {
			log.Fatalln(err)
			return
		}

		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, []byte(cemResult), "", "  ")
		if error != nil {
			log.Fatal("JSON parse error: ", error)
		}
		fmt.Println(prettyJSON.String())

	},
}

var cemGetRecommendationsCmd = &cobra.Command{
	Use:   "get-recommendations",
	Short: "Get Entity Recommendations",
	Long: `Retrieve recommendations for an entity

	Example Usage:
	$ cybr cem get-recommendations -p PLATFORM -a ACCOUNT_ID -e ENTITY_ID`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := GetToken()
		if err != nil {
			log.Fatalf("Failed to retrieve token file at %s. %s\n", CemSessionTokenPath, err)
		}

		cemQuery := &cem.EntityQuery{
			Platform:  CemPlatform,
			AccountId: CemAccountID,
			EntityId:  CemEntityID,
		}

		cemResult, err := cem.GetEntityRecommendations(token, cemQuery)
		if err != nil {
			log.Fatalln(err)
			return
		}

		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, []byte(cemResult), "", "  ")
		if error != nil {
			log.Fatal("JSON parse error: ", error)
		}
		fmt.Println(prettyJSON.String())

	},
}

var cemGetEntityDetailCmd = &cobra.Command{
	Use:   "get-entity-detail",
	Short: "Get Entity Details",
	Long: `Retrieve the details of a specific entity within a platform and workspace

	Example Usage:
	$ cybr cem get-entity-detail -p PLATFORM -a ACCOUNT_ID -e ENTITY_ID`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := GetToken()
		if err != nil {
			log.Fatalf("Failed to retrieve token file at %s. %s\n", CemSessionTokenPath, err)
		}

		cemQuery := &cem.EntityQuery{
			Platform:  CemPlatform,
			AccountId: CemAccountID,
			EntityId:  CemEntityID,
		}

		cemResult, err := cem.GetEntityDetail(token, cemQuery)
		if err != nil {
			log.Fatalln(err)
			return
		}

		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, []byte(cemResult), "", "  ")
		if error != nil {
			log.Fatal("JSON parse error: ", error)
		}
		fmt.Println(prettyJSON.String())

	},
}

var cemGetEntitiesCmd = &cobra.Command{
	Use:   "get-entities",
	Short: "Get Entities",
	Long: `Search for entities on any platform and retrieve entity details.
	
	Example Usage:
	$ cybr cem get-entities -p PLATFORM -a ACCOUNT_ID`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := GetToken()
		if err != nil {
			log.Fatalf("Failed to retrieve token file at %s. %s\n", CemSessionTokenPath, err)
		}

		if CemNonFullAdmin && CemFullAdmin {
			log.Fatal("Cannot set both --non-full-admin and --full-admin at the same time")
		}

		if CemNonShadowAdmin && CemShadowAdmin {
			log.Fatal("Cannot set both --non-shadow-admin and --shadow-admin at the same time")
		}

		inputFullAdmin := ""
		inputShadowAdmin := ""

		if CemFullAdmin {
			inputFullAdmin = "true"
		}
		if CemNonFullAdmin {
			inputFullAdmin = "false"
		}
		if CemShadowAdmin {
			inputShadowAdmin = "true"
		}
		if CemNonShadowAdmin {
			inputShadowAdmin = "false"
		}

		cemQuery := &cem.GetEntitiesQuery{
			Platform:    CemPlatform,
			AccountId:   CemAccountID,
			FullAdmin:   inputFullAdmin,
			ShadowAdmin: inputShadowAdmin,
			NextToken:   CemNextToken,
		}

		cemResult, err := cem.GetEntities(token, cemQuery)
		if err != nil {
			log.Fatalln(err)
			return
		}

		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, []byte(cemResult), "", "  ")
		if error != nil {
			log.Fatal("JSON parse error: ", error)
		}
		fmt.Println(prettyJSON.String())
	},
}

func init() {

	// Login
	cemLoginCmd.Flags().StringVarP(&CemOrganization, "organization", "o", "", "Username to logon PAS REST API using")
	cemLoginCmd.MarkFlagRequired("organization")
	cemLoginCmd.Flags().BoolVar(&CemNonInteractive, "non-interactive", false, "If detected, will retrieve the API key from the CEM_APIKEY environment variable")

	// Get Entities
	cemGetEntitiesCmd.Flags().StringVarP(&CemPlatform, "platform", "p", "", "Platform Name")
	cemGetEntitiesCmd.Flags().StringVarP(&CemAccountID, "account-id", "a", "", "Account ID")
	cemGetEntitiesCmd.Flags().BoolVar(&CemNonFullAdmin, "non-full-admin", false, "Get non-full admin entities only. Cannot be used with --full-admin")
	cemGetEntitiesCmd.Flags().BoolVar(&CemNonShadowAdmin, "non-shadow-admin", false, "Get non-shadow admin entities only.  Cannot be used with --shadow-admin")
	cemGetEntitiesCmd.Flags().BoolVar(&CemFullAdmin, "full-admin", false, "Get full admin entities only.  Cannot be used with --non-full-admin")
	cemGetEntitiesCmd.Flags().BoolVar(&CemShadowAdmin, "shadow-admin", false, "Get shadow admin entities only.  Cannot be used with --non-shadow-admin")
	cemGetEntitiesCmd.Flags().StringVarP(&CemNextToken, "next-token", "n", "", "The token for paging the entities.")

	// Get Entity Detail
	cemGetEntityDetailCmd.Flags().StringVarP(&CemPlatform, "platform", "p", "", "Platform Name")
	cemGetEntityDetailCmd.MarkFlagRequired("platform")
	cemGetEntityDetailCmd.Flags().StringVarP(&CemAccountID, "account-id", "a", "", "Account ID")
	cemGetEntityDetailCmd.MarkFlagRequired("account-id")
	cemGetEntityDetailCmd.Flags().StringVarP(&CemEntityID, "entity-id", "e", "", "Entity ID")
	cemGetEntityDetailCmd.MarkFlagRequired("entity-id")

	// Get Entity Recommendations
	cemGetRecommendationsCmd.Flags().StringVarP(&CemPlatform, "platform", "p", "", "Platform Name")
	cemGetRecommendationsCmd.MarkFlagRequired("platform")
	cemGetRecommendationsCmd.Flags().StringVarP(&CemAccountID, "account-id", "a", "", "Account ID")
	cemGetRecommendationsCmd.MarkFlagRequired("account-id")
	cemGetRecommendationsCmd.Flags().StringVarP(&CemEntityID, "entity-id", "e", "", "Entity ID")
	cemGetRecommendationsCmd.MarkFlagRequired("entity-id")

	// Get Entity Remediations
	cemGetRemediationsCmd.Flags().StringVarP(&CemPlatform, "platform", "p", "", "Platform Name")
	cemGetRemediationsCmd.MarkFlagRequired("platform")
	cemGetRemediationsCmd.Flags().StringVarP(&CemAccountID, "account-id", "a", "", "Account ID")
	cemGetRemediationsCmd.MarkFlagRequired("account-id")
	cemGetRemediationsCmd.Flags().StringVarP(&CemEntityID, "entity-id", "e", "", "Entity ID")
	cemGetRemediationsCmd.MarkFlagRequired("entity-id")

	// Add the sub-commands to "cem" command
	cemCmd.AddCommand(cemLoginCmd)
	cemCmd.AddCommand(cemGetEntitiesCmd)
	cemCmd.AddCommand(cemGetEntityDetailCmd)
	cemCmd.AddCommand(cemGetRecommendationsCmd)
	cemCmd.AddCommand(cemGetRemediationsCmd)
	cemCmd.AddCommand(cemGetAccountsCmd)

	rootCmd.AddCommand(cemCmd)
}
