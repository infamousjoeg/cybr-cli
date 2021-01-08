package cmd

import (
	"fmt"
	"log"
	"strings"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
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
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Account actions for PAS REST API",
	Long: `All account actions that can be taken via PAS REST API.
	
	Example Usage:
	List all accounts: $ cybr accounts list
	Get a Account details: $ cybr accounts get 234_1`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		apps, err := client.ListAccounts(&pasapi.ListAccountQueryParams{})
		if err != nil {
			log.Fatalf("Failed to retrieve a list of all accounts. %s", err)
			return
		}

		prettyprint.PrintJSON(apps)
	},
}

var listAccountsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long: `List all accounts the logged on user can read from PAS REST API.
	
	Example Usage:
	$ cybr accounts list`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		query := &pasapi.ListAccountQueryParams{
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
	$ cybr accounts get 24_1`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfig()
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

// The content will look like
// port=something, sp
func platformPropertiesToMap(content string) (map[string]string, error) {
	if content == "" {
		return nil, nil
	}

	if !strings.Contains(content, "=") {
		return nil, fmt.Errorf("Invalid platform prop content. The provided content does not container a '='")
	}

	m := make(map[string]string)

	// TODO: Gotta be a better way to do this
	replaceWith := "^||||^"

	// If the address or property contains a `\,` then replace
	content = strings.ReplaceAll(content, "\\,", replaceWith)
	props := strings.Split(content, ",")
	for _, prop := range props {
		if !strings.Contains(prop, "=") {
			return nil, fmt.Errorf("Property '%s' is invalid because it does not contain a '=' to seperate key from value", prop)
		}
		kvs := strings.SplitN(prop, "=", 2)
		key := strings.Trim(kvs[0], " ")
		value := strings.Trim(strings.ReplaceAll(kvs[1], replaceWith, ","), " ")
		m[key] = value
	}

	return m, nil
}

var addAccountsCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an account",
	Long: `Add an account to PAS.
	
	Example Usage:
	$ cybr accounts add -s SafeName -p platformID -u username -a 10.0.0.1 -t password -s SuperSecret`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		platformProps, err := platformPropertiesToMap(PlatformProperties)
		if err != nil {
			log.Fatalf("Failed to parse platform properties. %s", err)
		}

		newAccount := pasapi.AddAccountRequest{
			Name:       Name,
			Address:    Address,
			UserName:   Username,
			PlatformID: PlatformID,
			SafeName:   Safe,
			SecretType: SecretType,
			Secret:     Secret,
			SecretManagement: pasapi.SecretManagement{
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
	$ cybr accounts delete 24_1`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfig()
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

	// Add cmd to account cmd
	accountsCmd.AddCommand(listAccountsCmd)
	accountsCmd.AddCommand(getAccountsCmd)
	accountsCmd.AddCommand(addAccountsCmd)
	accountsCmd.AddCommand(deleteAccountsCmd)

	// Add accounts cmd to root
	rootCmd.AddCommand(accountsCmd)
}
