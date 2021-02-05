package cmd

import (
	"fmt"
	"log"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/prettyprint"
	"github.com/spf13/cobra"
)

// Safe is the safe name to filter on
var Safe string

var (
	// SafeName is the name of the safe to create
	SafeName string
	// Description is the description for the safe to create
	Description string
	// OLACEnabled is the boolean value of whether object-level access is enabled
	OLACEnabled bool
	// ManagingCPM is the name of the CPM User that manages accounts in the safe
	ManagingCPM string
	// NumberOfVersionsRetention is the number of password versions to retain for accounts within
	NumberOfVersionsRetention int
	// NumberOfDaysRetention is the number of days to retain older password versions for
	NumberOfDaysRetention int
	// AutoPurgeEnabled is a boolean value as to whether to remove non-compliant accounts automatically
	AutoPurgeEnabled bool
	// SafeLocation is the location the safe will be created in the Secure Digital Vault (default: \\)
	SafeLocation string
	// TargetSafeName is used by the Update Safe endpoint to refer to
	TargetSafeName string
	// UseAccounts use account inside of safe
	UseAccounts bool
	// RetrieveAccounts retrieve accounts inside of safe
	RetrieveAccounts bool
	// ListAccounts list accounts inside of safe
	ListAccounts bool
	// AddAccounts add account inside of safe
	AddAccounts bool
	// UpdateAccountContent update account content inside of safe
	UpdateAccountContent bool
	// UpdateAccountProperties update account properties inside of safe
	UpdateAccountProperties bool
	// InitiateCPMAccountManagementOperations init a cpm account action in safe
	InitiateCPMAccountManagementOperations bool
	// SpecifyNextAccountContent specify next account content in safe
	SpecifyNextAccountContent bool
	// ManageSafe manage this safe
	ManageSafe bool
	// ManageSafeMembers manage members of this safe
	ManageSafeMembers bool
	// BackupSafe backup the safe
	BackupSafe bool
	// ViewAuditLog view audit logs of this safe
	ViewAuditLog bool
	// ViewSafeMembers view member so this safe
	ViewSafeMembers bool
	// AccessWithoutConfirmation access safe without confirmation
	AccessWithoutConfirmation bool
	// CreateFolders create folders in safe
	CreateFolders bool
	// DeleteFolders delete folders in safe
	DeleteFolders bool
	// MoveAccountsAndFolders move accounts and folders
	MoveAccountsAndFolders bool
	// MemberName name of the member being added to a safe
	MemberName string
	//SearchIn search in Vault or Domain
	SearchIn string
	// MembershipExpirationDate when membership will expire
	MembershipExpirationDate string
)

var safesCmd = &cobra.Command{
	Use:   "safes",
	Short: "Safe actions for PAS REST API",
	Long: `All safe actions that can be taken via PAS REST API.
	
	Example Usage:
	List All Safes: $ cybr safes list
	List Safe Members: $ cybr safes member list -s SafeName
	Add Safe: $ cybr safes add -s SafeName -d Description --cpm ManagingCPM --days 0`,
	Aliases: []string{"safe"},
}

var listSafesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all safes",
	Long: `List all safes the logged on user can read from PAS REST API.
	
	Example Usage:
	$ cybr safes list`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}
		// List All Safes
		safes, err := client.ListSafes()
		if err != nil {
			log.Fatalf("Failed to retrieve a list of all safes. %s", err)
			return
		}
		// Pretty print returned object as JSON blob
		prettyprint.PrintJSON(safes)
	},
}

var listMembersCmd = &cobra.Command{
	Use:   "list-members",
	Short: "List all safe members on safes or specific safe",
	Long: `List all safe members on safes or a specific safe that
	the user logged on can read from PAS REST API.
	
	Example Usage:
	$ cybr safes list-members -s SafeName`,
	Aliases: []string{"list-member"},
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}
		// Add a safe with the configuration options given via CLI subcommands
		members, err := client.ListSafeMembers(Safe)
		if err != nil {
			log.Fatalf("Failed to retrieve a list of all safe members for %s. %s", Safe, err)
			return
		}
		// Pretty print returned object as JSON blob
		prettyprint.PrintJSON(members)
	},
}

var addMembersCmd = &cobra.Command{
	Use:   "add-member",
	Short: "Add a member to a safe with specific permissions",
	Long: `This method adds an existing user as a Safe member.
	The user who runs this web service requires Manage Safe Members permissions in the Vault.
	
	Example Usage:
	$ cybr safes add-member -s SafeName -m MemberName --retrieve-account`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		newMember := requests.AddSafeMember{
			Member: requests.AddSafeMemberInternal{
				MemberName:               MemberName,
				SearchIn:                 SearchIn,
				MembershipExpirationDate: MembershipExpirationDate,
				Permissions: []requests.PermissionKeyValue{
					requests.PermissionKeyValue{
						Key:   "UseAccounts",
						Value: UseAccounts,
					},
					requests.PermissionKeyValue{
						Key:   "RetrieveAccounts",
						Value: RetrieveAccounts,
					},
					requests.PermissionKeyValue{
						Key:   "ListAccounts",
						Value: ListAccounts,
					},
					requests.PermissionKeyValue{
						Key:   "AddAccounts",
						Value: AddAccounts,
					},
					requests.PermissionKeyValue{
						Key:   "UpdateAccountContent",
						Value: UpdateAccountContent,
					},
					requests.PermissionKeyValue{
						Key:   "UpdateAccountProperties",
						Value: UpdateAccountProperties,
					},
					requests.PermissionKeyValue{
						Key:   "InitiateCPMAccountManagementOperations",
						Value: InitiateCPMAccountManagementOperations,
					},
					requests.PermissionKeyValue{
						Key:   "SpecifyNextAccountContent",
						Value: SpecifyNextAccountContent,
					},
					requests.PermissionKeyValue{
						Key:   "ManageSafe",
						Value: ManageSafe,
					},
					requests.PermissionKeyValue{
						Key:   "ManageSafeMembers",
						Value: ManageSafeMembers,
					},
					requests.PermissionKeyValue{
						Key:   "BackupSafe",
						Value: BackupSafe,
					},
					requests.PermissionKeyValue{
						Key:   "ViewAuditLog",
						Value: ViewAuditLog,
					},
					requests.PermissionKeyValue{
						Key:   "ViewSafeMembers",
						Value: ViewSafeMembers,
					},
					requests.PermissionKeyValue{
						Key:   "AccessWithoutConfirmation",
						Value: AccessWithoutConfirmation,
					},
					requests.PermissionKeyValue{
						Key:   "CreateFolders",
						Value: CreateFolders,
					},
					requests.PermissionKeyValue{
						Key:   "DeleteFolders",
						Value: DeleteFolders,
					},
					requests.PermissionKeyValue{
						Key:   "MoveAccountsAndFolders",
						Value: MoveAccountsAndFolders,
					},
				},
			},
		}

		// Add a safe with the configuration options given via CLI subcommands
		err = client.AddSafeMember(Safe, newMember)
		if err != nil {
			log.Fatalf("Failed to add member '%s' to safe '%s'. %s", MemberName, Safe, err)
			return
		}
		fmt.Printf("Successfully added member '%s' to safe '%s'\n", MemberName, Safe)
	},
}

var removeMembersCmd = &cobra.Command{
	Use:   "remove-member",
	Short: "Remove a member from a safe",
	Long: `This method removes a specific member from a Safe.
	The user who runs this web service requires Manage Safe Members permissions in the Vault.
	
	Example Usage:
	$ cybr safes remove-member -s SafeName -m MemberName`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		// Add a safe with the configuration options given via CLI subcommands
		err = client.RemoveSafeMember(Safe, MemberName)
		if err != nil {
			log.Fatalf("Failed to add member '%s' to safe '%s'. %s", MemberName, Safe, err)
			return
		}

		fmt.Printf("Successfully removed member '%s' from safe '%s'\n", MemberName, Safe)
	},
}

var addSafeCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a safe",
	Long: `Add a safe and configure its retention and location
	via the PAS REST API.
	
	Example Usage:
	$ cybr safes add -s SafeName -d Description --cpm ManagingCPM --days 0`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}
		// Build body of the request
		body := requests.AddSafe{
			SafeName:              SafeName,
			Description:           Description,
			OLACEnabled:           OLACEnabled,
			ManagingCPM:           ManagingCPM,
			NumberOfDaysRetention: NumberOfDaysRetention,
			AutoPurgeEnabled:      AutoPurgeEnabled,
			SafeLocation:          SafeLocation,
		}
		// Add the safe with config declared above
		err = client.AddSafe(body)
		if err != nil {
			log.Fatalf("Failed to add the safe named %s. %s", SafeName, err)
			return
		}
		fmt.Printf("Successfully added safe %s.\n", SafeName)
	},
}

var deleteSafeCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a safe",
	Long: `Delete a safe via the PAS REST API. If the safe has a retention policy
	of days that is greater than 0, the safe will be marked for deletion until
	the amount of days assigned are met.
	
	Example Usage:
	$ cybr safes delete -s SafeName`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}
		// Delete the safe
		err = client.DeleteSafe(SafeName)
		if err != nil {
			log.Fatalf("Failed to delete the safe named %s. %s", SafeName, err)
			return
		}

		fmt.Printf("Successfully deleted safe %s.\n", SafeName)
	},
}

var updateSafeCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a safe",
	Long: `Update a safe via the PAS REST API. Only the options provided will be modified.
	
	Example Usage:
	$ cybr safes update -t TargetSafeName -s NewSafeName -d NewDesc`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}
		// Build body of the request
		body := requests.UpdateSafe{
			SafeName:    SafeName,
			Description: Description,
			OLACEnabled: OLACEnabled,
			ManagingCPM: ManagingCPM,
		}
		// Update the safe
		response, err := client.UpdateSafe(TargetSafeName, body)
		if err != nil {
			log.Fatalf("Failed to update the safe named %s. %s", TargetSafeName, err)
			return
		}
		// Pretty print returned object as JSON blob
		prettyprint.PrintJSON(response)
	},
}

func init() {
	listMembersCmd.Flags().StringVarP(&Safe, "safe", "s", "", "Safe name to filter request on")
	listMembersCmd.MarkFlagRequired("safe")

	addSafeCmd.Flags().StringVarP(&SafeName, "safe", "s", "", "Safe name to create")
	addSafeCmd.Flags().StringVarP(&Description, "desc", "d", "", "Description of the safe created")
	addSafeCmd.Flags().BoolVarP(&OLACEnabled, "olac", "O", false, "Enable object-level access control (OLAC) on safe (cannot be reversed)")
	addSafeCmd.Flags().StringVarP(&ManagingCPM, "cpm", "", "PasswordManager", "Set the Managing CPM user to something other than PasswordManager")
	addSafeCmd.Flags().IntVarP(&NumberOfDaysRetention, "days", "", 7, "Number of days to retain password versions for")
	addSafeCmd.Flags().BoolVarP(&AutoPurgeEnabled, "auto-purge", "P", false, "Whether to automatically purge accounts after a number of records is met")
	addSafeCmd.Flags().StringVarP(&SafeLocation, "location", "l", "\\", "The location of the Safe in the Secure Digital Vault")
	addSafeCmd.MarkFlagRequired("safe")
	addSafeCmd.MarkFlagRequired("desc")

	deleteSafeCmd.Flags().StringVarP(&SafeName, "safe", "s", "", "Safe name to delete")
	deleteSafeCmd.MarkFlagRequired("safe")

	updateSafeCmd.Flags().StringVarP(&TargetSafeName, "target-safe", "t", "", "Safe name to update")
	updateSafeCmd.Flags().StringVarP(&SafeName, "safe", "s", "", "New safe name to change to")
	updateSafeCmd.Flags().StringVarP(&Description, "desc", "d", "", "New description to change to")
	updateSafeCmd.Flags().BoolVarP(&OLACEnabled, "olac", "O", false, "Enable object-level access control (OLAC) on safe (cannot be disabled)")
	updateSafeCmd.Flags().StringVarP(&ManagingCPM, "cpm", "", "", "New managing CPM user to change to")
	updateSafeCmd.MarkFlagRequired("target-safe")

	// add-member
	addMembersCmd.Flags().StringVarP(&Safe, "safe", "s", "", "Name of the safe")
	addMembersCmd.MarkFlagRequired("safe")
	addMembersCmd.Flags().StringVarP(&MemberName, "member-name", "m", "", "Name of member being added to the desired safe")
	addMembersCmd.MarkFlagRequired("member-name")
	addMembersCmd.Flags().StringVarP(&MemberName, "search-in", "i", "Vault", "Search in Domain or Vault")
	addMembersCmd.Flags().StringVarP(&MembershipExpirationDate, "member-expiration-date", "e", "", "When the membership will expire")
	addMembersCmd.Flags().BoolVar(&UseAccounts, "use-accounts", false, "Use accounts in safe")
	addMembersCmd.Flags().BoolVar(&RetrieveAccounts, "retrieve-accounts", false, "Retrieve accounts in safe")
	addMembersCmd.Flags().BoolVar(&ListAccounts, "list-accounts", false, "List accounts in safe")
	addMembersCmd.Flags().BoolVar(&AddAccounts, "add-accounts", false, "Add accounts to safe")
	addMembersCmd.Flags().BoolVar(&UpdateAccountContent, "update-account-content", false, "Update account content in safe")
	addMembersCmd.Flags().BoolVar(&UpdateAccountProperties, "update-account-properties", false, "Update account properties in safe")
	addMembersCmd.Flags().BoolVar(&InitiateCPMAccountManagementOperations, "init-cpm-account-managment-operations", false, "Perform cpm actions on accounts inside of safe")
	addMembersCmd.Flags().BoolVar(&SpecifyNextAccountContent, "specify-next-account-content", false, "Specify next account's content within safe")
	addMembersCmd.Flags().BoolVar(&ManageSafe, "manage-safe", false, "Manage the safe")
	addMembersCmd.Flags().BoolVar(&ManageSafeMembers, "manage-safe-members", false, "Manage members of the safe")
	addMembersCmd.Flags().BoolVar(&BackupSafe, "backup-safe", false, "Backup the safe")
	addMembersCmd.Flags().BoolVar(&ViewAuditLog, "view-audit-log", false, "View audit log of safe")
	addMembersCmd.Flags().BoolVar(&ViewSafeMembers, "view-safe-members", false, "View the safe members")
	addMembersCmd.Flags().BoolVar(&AccessWithoutConfirmation, "access-content-without-confirmation", false, "Access safe content without confirmation")
	addMembersCmd.Flags().BoolVar(&CreateFolders, "create-folders", false, "Create folders within safe")
	addMembersCmd.Flags().BoolVar(&DeleteFolders, "delete-folders", false, "Delete folders within safe")
	addMembersCmd.Flags().BoolVar(&MoveAccountsAndFolders, "move-accounts-and-folders", false, "Move accounts and folders")

	// remove-member
	removeMembersCmd.Flags().StringVarP(&Safe, "safe", "s", "", "Name of the safe")
	removeMembersCmd.MarkFlagRequired("safe-name")
	removeMembersCmd.Flags().StringVarP(&MemberName, "member-name", "m", "", "Name of member being removed from the safe")
	removeMembersCmd.MarkFlagRequired("member-name")

	safesCmd.AddCommand(listSafesCmd)
	safesCmd.AddCommand(listMembersCmd)
	safesCmd.AddCommand(addSafeCmd)
	safesCmd.AddCommand(deleteSafeCmd)
	safesCmd.AddCommand(updateSafeCmd)
	safesCmd.AddCommand(addMembersCmd)
	safesCmd.AddCommand(removeMembersCmd)
	rootCmd.AddCommand(safesCmd)
}
