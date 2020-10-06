package cmd

import (
	"log"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/prettyprint"
	"github.com/spf13/cobra"
)

var (
	// AppID is the application identity to filter on
	AppID string
	// Location is the folder location the Application is located in
	Location string
)

var applicationsCmd = &cobra.Command{
	Use:   "applications",
	Short: "Applications actions for PAS REST API",
	Long: `All applications actions that can be taken via PAS REST API.
	
	Example Usage:
	List All Applications at Root: $ cybr applications list
	List All Applications at \Applications: $ cybr applications list -l \\Applications
	List All Authentication Methods: $ cybr applications methods list -a AppID`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}
		// List All Safes
		apps, err := client.ListApplications(Location)
		if err != nil {
			log.Fatalf("Failed to retrieve a list of all applications. %s", err)
			return
		}
		// Pretty print returned object as JSON blob
		prettyprint.PrintJSON(apps)
	},
}

var listApplicationsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all applications",
	Long: `List all applications the logged on user can read from PAS REST API.
	
	Example Usage:
	$ cybr applications list`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}
		// List All Safes
		apps, err := client.ListApplications(Location)
		if err != nil {
			log.Fatalf("Failed to retrieve a list of all applications. %s", err)
			return
		}
		// Pretty print returned object as JSON blob
		prettyprint.PrintJSON(apps)
	},
}

var listMethodsCmd = &cobra.Command{
	Use:   "methods list",
	Short: "List all authn methods on a specific application",
	Long: `List all authentication methods on a specific application
	that the user logged on can read from PAS REST API.
	
	Example Usage:
	$ cybr applications methods list -a AppID`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}
		// List all Safe Members for specific safe ""
		methods, err := client.ListApplicationAuthenticationMethods(AppID)
		if err != nil {
			log.Fatalf("Failed to retrieve a list of all application methods for %s. %s", Safe, err)
			return
		}
		// Pretty print returned object as JSON blob
		prettyprint.PrintJSON(methods)
	},
}

func init() {
	listApplicationsCmd.Flags().StringVarP(&Location, "location", "l", "\\", "Location of the application in EPV")
	listMethodsCmd.Flags().StringVarP(&AppID, "app-id", "a", "", "Application identity to filter request on")
	listMethodsCmd.MarkFlagRequired("app-id")
	applicationsCmd.Flags().StringVarP(&Location, "location", "l", "\\", "Location of the application in EPV")
	applicationsCmd.AddCommand(listApplicationsCmd)
	applicationsCmd.AddCommand(listMethodsCmd)
	rootCmd.AddCommand(applicationsCmd)
}
