package cmd

import (
	"fmt"
	"log"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/queries"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/shared"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/prettyprint"
	"github.com/spf13/cobra"
)

var (
	// AccountID is a specific account ID used for updating, getting and deleting
	AccountID string

	// PlatformID of the account object being added
	PlatformID string

	// Name name of the account object
	Name string

	// Address of the account
	Address string

	// SecretType of the account
	SecretType string

	// Secret of the account
	Secret string

	// AutomaticManagementEnabled if account will be managed
	AutomaticManagementEnabled bool

	// ManualManagementReason reason account is not being managed
	ManualManagementReason string

	// PlatformProperties for account
	PlatformProperties string

	// Search List of keywords to search for in accounts, separated by a space.
	Search string

	// SearchType Get accounts that either contain or start with the value specified in the Search parameter. Valid values: contains (default) or startswith
	SearchType string

	// Sort Property or properties by which to sort returned accounts, followed by asc (default) or desc to control sort direction. Separate multiple properties with commas, up to a maximum of three properties.
	Sort string

	// Offset of the first account that is returned in the collection of results.
	Offset int

	// Limit Maximum number of returned accounts. If not specified, the default value is 50. The maximum number that can be specified is 1000.
	Limit int

	// Filter Search for accounts filtered by safeName or modificationTime
	Filter string

	// Reason to access account
	Reason string

	// TicketingSystemName name of the ticketing system
	TicketingSystemName string

	// TicketID the ticket ID
	TicketID string

	// Version of the secret/password being retrieved
	Version int

	// ChangeEntireGroup change account group
	ChangeEntireGroup bool
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Account actions for PAS REST API",
	Long: `All account actions that can be taken via PAS REST API.
	
	Example Usage:
	List all accounts: $ cybr accounts list
	Get a Account details: $ cybr accounts get 234_1`,
	Aliases: []string{"account"},
}

var listAccountsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long: `List all accounts the logged on user can read from PAS REST API.
	
	Example Usage:
	$ cybr accounts list`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		query := &queries.ListAccounts{
			Search:     Search,
			SearchType: SearchType,
			Sort:       Sort,
			Offset:     Offset,
			Limit:      Limit,
			Filter:     Filter,
		}

		apps, err := client.ListAccounts(query)
		if err != nil {
			log.Fatalf("Failed to retrieve a list of all accounts. %s", err)
			return
		}

		prettyprint.PrintJSON(apps)
	},
}

var getAccountsCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific account",
	Long: `Get a specific account from PAS REST API.
	
	Example Usage:
	$ cybr accounts get -i 24_1`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		apps, err := client.GetAccount(AccountID)
		if err != nil {
			log.Fatalf("Failed to retrieve account '%s'. %s", AccountID, err)
			return
		}

		prettyprint.PrintJSON(apps)
	},
}

var addAccountsCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an account",
	Long: `Add an account to PAS.
	
	Example Usage:
	$ cybr accounts add -s SafeName -p platformID -u username -a 10.0.0.1 -t password -s SuperSecret`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		platformProps, err := keyValueStringToMap(PlatformProperties)
		if err != nil {
			log.Fatalf("Failed to parse platform properties. %s", err)
		}

		newAccount := requests.AddAccount{
			Name:       Name,
			Address:    Address,
			UserName:   Username,
			PlatformID: PlatformID,
			SafeName:   Safe,
			SecretType: SecretType,
			Secret:     Secret,
			SecretManagement: shared.SecretManagement{
				AutomaticManagementEnabled: AutomaticManagementEnabled,
				ManualManagementReason:     ManualManagementReason,
			},
			PlatformAccountProperties: platformProps,
		}

		apps, err := client.AddAccount(newAccount)
		if err != nil {
			log.Fatalf("Failed to add account. %s", err)
			return
		}

		prettyprint.PrintJSON(apps)
	},
}

var deleteAccountsCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a specific account",
	Long: `Delete a specific account from PAS REST API.
	
	Example Usage:
	$ cybr accounts delete -i 24_1`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		err = client.DeleteAccount(AccountID)
		if err != nil {
			log.Fatalf("Failed to delete account '%s'. %s", AccountID, err)
			return
		}

		fmt.Printf("Successfully deleted account with id '%s'\n", AccountID)
	},
}

var getPasswordAccountCmd = &cobra.Command{
	Use:   "get-password",
	Short: "Get password of a specific account",
	Long: `This method enables users to retrieve the password or SSH key of an existing account that is identified by its Account ID. It enables users to specify a reason and ticket ID, if required.
	
	Example Usage:
	$ cybr accounts get-password -i 24_1`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		request := requests.GetAccountPassword{
			Reason:              Reason,
			TicketingSystemName: TicketingSystemName,
			TicketID:            TicketID,
			Version:             Version,
		}

		response, err := client.GetAccountPassword(AccountID, request)
		if err != nil {
			log.Fatalf("%s", err)
			return
		}

		fmt.Println(response)
	},
}

var verifyAccountCmd = &cobra.Command{
	Use:   "verify",
	Short: "Mark an account for verification",
	Long: `This method marks an account for credential verification
	
	Example Usage:
	$ cybr accounts verify -i 24_1`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		err = client.VerifyAccountCredentials(AccountID)
		if err != nil {
			log.Fatalf("%s", err)
			return
		}

		fmt.Printf("Successfully marked account '%s' for verification\n", AccountID)
	},
}

var changeAccountCmd = &cobra.Command{
	Use:   "change",
	Short: "Mark an account for change",
	Long: `This method marks an account for credential change
	
	Example Usage:
	$ cybr accounts change -i 24_1`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		err = client.ChangeAccountCredentials(AccountID, ChangeEntireGroup)
		if err != nil {
			log.Fatalf("%s", err)
			return
		}

		fmt.Printf("Successfully marked account '%s' for change\n", AccountID)
	},
}

var reconcileAccountCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Mark an account for reconciliation",
	Long: `This method marks an account for credential reconciliation
	
	Example Usage:
	$ cybr accounts reconcile -i 24_1`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		err = client.ReconileAccountCredentials(AccountID)
		if err != nil {
			log.Fatalf("%s", err)
			return
		}

		fmt.Printf("Successfully marked account '%s' for reconciliation\n", AccountID)
	},
}

var moveAccountCmd = &cobra.Command{
	Use:   "move",
	Short: "Move an account to a different safe",
	Long: `Move an account to a different safe

	Example Usage:
	$ cybr accounts move -i 24_1 -s newSafeName`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		account, err := client.GetAccount(AccountID)
		if err != nil {
			log.Fatalf("%s", err)
		}

		secret, err := client.GetAccountPassword(AccountID, requests.GetAccountPassword{})
		if err != nil {
			log.Fatalf("%s", err)
		}

		newAccount := requests.AddAccount{
			Name:                      account.Name,
			Address:                   account.Address,
			UserName:                  account.UserName,
			PlatformID:                account.PlatformID,
			SafeName:                  Safe,
			SecretType:                account.SecretType,
			Secret:                    secret,
			PlatformAccountProperties: account.PlatformAccountProperties,
			SecretManagement:          account.SecretManagement,
		}

		createdAccount, err := client.AddAccount(newAccount)
		if err != nil {
			log.Fatalf("%s", err)
		}

		err = client.DeleteAccount(AccountID)
		if err != nil {
			log.Fatalf("%s", err)
			return
		}

		prettyprint.PrintJSON(createdAccount)
	},
}

func init() {
	// Listing an account
	listAccountsCmd.Flags().StringVarP(&Search, "search", "s", "", "List of keywords to search for in accounts, separated by a space")
	listAccountsCmd.Flags().StringVarP(&SearchType, "search-type", "t", "", "Get accounts that either contain or start with the value specified in the Search parameter. Valid values: contains (default) or startswith")
	listAccountsCmd.Flags().StringVarP(&Sort, "sort", "r", "", "Property or properties by which to sort returned accounts, followed by asc (default) or desc to control sort direction. Separate multiple properties with commas, up to a maximum of three properties")
	listAccountsCmd.Flags().IntVarP(&Offset, "offset", "o", 0, "Offset of the first account that is returned in the collection of results")
	listAccountsCmd.Flags().IntVarP(&Limit, "limit", "l", 50, "Maximum number of returned accounts. If not specified, the default value is 50. The maximum number that can be specified is 1000")
	listAccountsCmd.Flags().StringVarP(&Filter, "filter", "f", "", "Search for accounts filtered by safeName or modificationTime")

	// Getting an account
	getAccountsCmd.Flags().StringVarP(&AccountID, "account-id", "i", "", "Account ID to list from")
	getAccountsCmd.MarkFlagRequired("account-id")

	// Creating an account
	addAccountsCmd.Flags().StringVarP(&Name, "name", "n", "", "The name of the account object being created. Will use auto-generated name if not provided")
	addAccountsCmd.Flags().StringVarP(&Address, "address", "a", "", "Address of the account object")
	addAccountsCmd.Flags().StringVarP(&Username, "username", "u", "", "Username of the account object")
	addAccountsCmd.Flags().StringVarP(&PlatformID, "platform-id", "p", "", "Platform ID of the account object")
	addAccountsCmd.Flags().StringVarP(&Safe, "safe", "s", "", "Safe name of the account object")
	addAccountsCmd.Flags().StringVarP(&SecretType, "secret-type", "t", "", "Secret type of the account object. e.g. password, accessKey, sshKey")
	addAccountsCmd.MarkFlagRequired("secret-type")
	addAccountsCmd.Flags().StringVarP(&Secret, "secret", "c", "", "Secret of the account object")
	addAccountsCmd.MarkFlagRequired("secret")
	addAccountsCmd.Flags().StringVarP(&PlatformProperties, "platform-properties", "e", "", "Extra platform properties. e.g. port=22,UseSudoOnReconcile=yes,CustomField=custom")
	addAccountsCmd.Flags().BoolVarP(&AutomaticManagementEnabled, "automatic-management", "m", false, "If set will automatically managed the onboarded account")
	addAccountsCmd.Flags().StringVarP(&ManualManagementReason, "manual-management-reason", "r", "", "The reason the account object is not being managed")

	// Delete an account
	deleteAccountsCmd.Flags().StringVarP(&AccountID, "account-id", "i", "", "Account ID to delete")
	deleteAccountsCmd.MarkFlagRequired("account-id")

	// Get password for account
	getPasswordAccountCmd.Flags().StringVarP(&AccountID, "account-id", "i", "", "Account ID to retrieve password value of")
	getPasswordAccountCmd.MarkFlagRequired("account-id")
	getPasswordAccountCmd.Flags().IntVarP(&Version, "version", "v", 0, "Version of the account password")
	getPasswordAccountCmd.Flags().StringVarP(&Reason, "reason", "r", "", "Reason for retriving account password")
	getPasswordAccountCmd.Flags().StringVarP(&TicketingSystemName, "ticketing-system", "s", "", "Ticketing system name")
	getPasswordAccountCmd.Flags().StringVarP(&TicketID, "ticket-id", "t", "", "The ticket ID related to the ticketing system")

	// verify account
	verifyAccountCmd.Flags().StringVarP(&AccountID, "account-id", "i", "", "Account ID to verify")
	verifyAccountCmd.MarkFlagRequired("account-id")

	// change account
	changeAccountCmd.Flags().StringVarP(&AccountID, "account-id", "i", "", "Account ID to change")
	changeAccountCmd.MarkFlagRequired("account-id")
	changeAccountCmd.Flags().BoolVarP(&ChangeEntireGroup, "change-entire-group", "c", false, "If account is part of account group, change the entire group")

	// reconcile
	reconcileAccountCmd.Flags().StringVarP(&AccountID, "account-id", "i", "", "Account ID to reconcile")
	reconcileAccountCmd.MarkFlagRequired("account-id")

	// move
	moveAccountCmd.Flags().StringVarP(&AccountID, "account-id", "i", "", "Account ID to move")
	moveAccountCmd.MarkFlagRequired("account-id")
	moveAccountCmd.Flags().StringVarP(&Safe, "safe", "s", "", "Safe name in which the account will be moved into")
	moveAccountCmd.MarkFlagRequired("safe")

	// Add cmd to account cmd
	accountsCmd.AddCommand(listAccountsCmd)
	accountsCmd.AddCommand(getAccountsCmd)
	accountsCmd.AddCommand(addAccountsCmd)
	accountsCmd.AddCommand(deleteAccountsCmd)
	accountsCmd.AddCommand(getPasswordAccountCmd)
	accountsCmd.AddCommand(verifyAccountCmd)
	accountsCmd.AddCommand(changeAccountCmd)
	accountsCmd.AddCommand(reconcileAccountCmd)
	accountsCmd.AddCommand(moveAccountCmd)

	// Add accounts cmd to root
	rootCmd.AddCommand(accountsCmd)
}
