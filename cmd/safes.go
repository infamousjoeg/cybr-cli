package cmd

import (
	"log"

	pasapi "github.com/infamousjoeg/pas-api-go/pkg/cybr/api"
	"github.com/infamousjoeg/pas-api-go/pkg/cybr/helpers/prettyprint"
	"github.com/spf13/cobra"
)

// Safe is the safe name to filter on
var Safe string

var safesCmd = &cobra.Command{
	Use:   "safes",
	Short: "Safe actions for PAS REST API",
	Long: `All safe actions that can be taken via PAS REST API.
	
	Example Usage:
	List All Safes: $ cybr safes list
	List Safe Members: $ cybr safes member list -s SafeName`,
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

var listCmd = &cobra.Command{
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
		// List all Safe Members for specific safe ""
		members, err := client.ListSafeMembers(Safe)
		if err != nil {
			log.Fatalf("Failed to retrieve a list of all safe members for %s. %s", Safe, err)
			return
		}
		// Pretty print returned object as JSON blob
		prettyprint.PrintJSON(members)
	},
}

func init() {
	listMembersCmd.Flags().StringVarP(&Safe, "safe", "s", "", "Safe name to filter request on")
	listMembersCmd.MarkFlagRequired("safe")
	safesCmd.AddCommand(listCmd)
	safesCmd.AddCommand(listMembersCmd)
	rootCmd.AddCommand(safesCmd)
}
