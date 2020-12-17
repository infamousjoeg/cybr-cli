package cmd

import (
	"fmt"
	"log"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
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
)

var safesCmd = &cobra.Command{
	Use:   "safes",
	Short: "Safe actions for PAS REST API",
	Long: `All safe actions that can be taken via PAS REST API.
	
	Example Usage:
	List All Safes: $ cybr safes list
	List Safe Members: $ cybr safes member list -s SafeName
	Add Safe: $ cybr safes add -s SafeName -d Description --cpm ManagingCPM --days 0`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfig()
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

var listSafesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all safes",
	Long: `List all safes the logged on user can read from PAS REST API.
	
	Example Usage:
	$ cybr safes list`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfig()
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
	Use:   "member list",
	Short: "List all safe members on safes or specific safe",
	Long: `List all safe members on safes or a specific safe that
	the user logged on can read from PAS REST API.
	
	Example Usage:
	$ cybr safes member list -s SafeName`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfig()
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

var addSafeCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a safe",
	Long: `Add a safe and configure its retention and location
	via the PAS REST API.
	
	Example Usage:
	$ cybr safes add -s SafeName -d Description --cpm ManagingCPM --days 0`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}
		// Build body of the request
		body := pasapi.AddSafeRequest{
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
		fmt.Printf("Successfully added safe %s.", SafeName)
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
		client, err := pasapi.GetConfig()
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

		fmt.Printf("Successfully deleted safe %s.", SafeName)
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
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}
		// Build body of the request
		body := pasapi.UpdateSafeRequest{
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

	safesCmd.AddCommand(listSafesCmd)
	safesCmd.AddCommand(listMembersCmd)
	safesCmd.AddCommand(addSafeCmd)
	safesCmd.AddCommand(deleteSafeCmd)
	safesCmd.AddCommand(updateSafeCmd)
	rootCmd.AddCommand(safesCmd)
}
