package cmd

import (
	"fmt"
	"log"
	"os"
	"syscall"

	local_cem "github.com/infamousjoeg/cybr-cli/pkg/cybr/cem"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/prettyprint"
	"github.com/quincycheng/cem-api-go/pkg/cem"
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
	Use:   "logon",
	Short: "Logon to CEM REST API",
	Long: `Authenticate to Cloud Entitlements Manager REST API.
	
	Example Usage:
	$ cybr cem login -a $ORGANIZATION
	
	For non-interactive logon:
	$ export ` + CemEnvAPIKey + `=<Your CEM access key>
	$ cybr cem login -a $ORGANIZATION --non-interactive`,
	Aliases: []string{"login"},
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

		err = local_cem.SaveToken(token, CemSessionTokenPath)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Successfully logged onto CEM (organization: %s)\n", CemOrganization)
	},
}

var cemGetAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Get Accounts",
	Long: `Retrieve workspaces grouped by platforms

	Example Usage:
	$ cybr cem get-accounts`,
	Aliases: []string{"get-accounts", "account"},
	Run: func(cmd *cobra.Command, args []string) {
		token, err := local_cem.GetToken(CemSessionTokenPath)
		if err != nil {
			log.Fatalf("Failed to retrieve token file at %s. %s\n", CemSessionTokenPath, err)
		}

		cemResult, err := cem.GetAccounts(token)
		if err != nil {
			log.Fatalln(err)
			return
		}

		prettyprint.PrintJSON(cemResult)
	},
}

var cemGetRemediationsCmd = &cobra.Command{
	Use:   "remediations",
	Short: "Get Entity Remediations",
	Long: `Retrieve Remediations for an entity

	Example Usage:
	$ cybr cem get-remediations -p PLATFORM -a ACCOUNT_ID -e ENTITY_ID`,
	Aliases: []string{"get-remediations", "remediation"},
	Run: func(cmd *cobra.Command, args []string) {
		token, err := local_cem.GetToken(CemSessionTokenPath)
		if err != nil {
			log.Fatalf("Failed to retrieve token file at %s. %s\n", CemSessionTokenPath, err)
		}

		cemQuery := &cem.EntityQuery{
			Platform:  CemPlatform,
			AccountID: CemAccountID,
			EntityID:  CemEntityID,
		}

		cemResult, err := cem.GetEntityRecommendations(token, cemQuery)
		if err != nil {
			log.Fatalln(err)
			return
		}

		prettyprint.PrintJSON(cemResult)

	},
}

var cemGetRecommendationsCmd = &cobra.Command{
	Use:   "recommendations",
	Short: "Get Entity Recommendations",
	Long: `Retrieve recommendations for an entity

	Example Usage:
	$ cybr cem get-recommendations -p PLATFORM -a ACCOUNT_ID -e ENTITY_ID`,
	Aliases: []string{"get-recommendations", "recommendation"},
	Run: func(cmd *cobra.Command, args []string) {
		token, err := local_cem.GetToken(CemSessionTokenPath)
		if err != nil {
			log.Fatalf("Failed to retrieve token file at %s. %s\n", CemSessionTokenPath, err)
		}

		cemQuery := &cem.EntityQuery{
			Platform:  CemPlatform,
			AccountID: CemAccountID,
			EntityID:  CemEntityID,
		}

		cemResult, err := cem.GetEntityRecommendations(token, cemQuery)
		if err != nil {
			log.Fatalln(err)
			return
		}

		prettyprint.PrintJSON(cemResult)

	},
}

var cemGetEntityDetailCmd = &cobra.Command{
	Use:   "entity-details",
	Short: "Get Entity Details",
	Long: `Retrieve the details of a specific entity within a platform and workspace

	Example Usage:
	$ cybr cem get-entity-detail -p PLATFORM -a ACCOUNT_ID -e ENTITY_ID`,
	Aliases: []string{"get-entity-details", "entity-detail"},
	Run: func(cmd *cobra.Command, args []string) {
		token, err := local_cem.GetToken(CemSessionTokenPath)
		if err != nil {
			log.Fatalf("Failed to retrieve token file at %s. %s\n", CemSessionTokenPath, err)
		}

		cemQuery := &cem.EntityQuery{
			Platform:  CemPlatform,
			AccountID: CemAccountID,
			EntityID:  CemEntityID,
		}

		cemResult, err := cem.GetEntityDetail(token, cemQuery)
		if err != nil {
			log.Fatalln(err)
			return
		}

		prettyprint.PrintJSON(cemResult)
	},
}

var cemGetEntitiesCmd = &cobra.Command{
	Use:   "entities",
	Short: "Get Entities",
	Long: `Search for entities on any platform and retrieve entity details.
	
	Example Usage:
	$ cybr cem get-entities -p PLATFORM -a ACCOUNT_ID`,
	Aliases: []string{"get-entities", "entity"},
	Run: func(cmd *cobra.Command, args []string) {
		token, err := local_cem.GetToken(CemSessionTokenPath)
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

		prettyprint.PrintJSON(cemResult)
	},
}

func init() {

	// Login
	cemLoginCmd.Flags().StringVarP(&CemOrganization, "organization", "o", "", "Organization identity for CEM")
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
