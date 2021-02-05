package cmd

import (
	"fmt"
	"log"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/prettyprint"
	"github.com/spf13/cobra"
)

var (
	// AppID is the application identity to filter on
	AppID string
	// Location is the folder location the Application is located in
	Location string
	// AuthType authentication method type
	AuthType string
	// AuthValue authentication method value
	AuthValue string
	// IsFolder used in path/hash authentication
	IsFolder bool
	// AllowInternalScripts allow internal script
	AllowInternalScripts bool
	// Desc app description
	Desc string
	// AccessPermittedFrom application access starting from
	AccessPermittedFrom int
	// AccessPermittedTo application access end at
	AccessPermittedTo int
	// ExpirationDate application expirey date
	ExpirationDate string
	// BusinessOwnerFName first name
	BusinessOwnerFName string
	// BusinessOwnerLName last name
	BusinessOwnerLName string
	// BusinessOwnerEmail email
	BusinessOwnerEmail string
	// BusinessOwnerPhone phone
	BusinessOwnerPhone string
	// Disabled application is disabled
	Disabled string
	// AppAuthnMethodID application authentication method ID
	AppAuthnMethodID string
)

var applicationsCmd = &cobra.Command{
	Use:   "applications",
	Short: "Applications actions for PAS REST API",
	Long: `All applications actions that can be taken via PAS REST API.
	
	Example Usage:
	List All Applications at Root: $ cybr applications list
	List All Applications at \Applications: $ cybr applications list -l \\Applications
	List All Authentication Methods: $ cybr applications methods list -a AppID`,
	Aliases: []string{"application", "app"},
}

var listApplicationsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all applications",
	Long: `List all applications the logged on user can read from PAS REST API.
	
	Example Usage:
	$ cybr applications list`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfigWithLogger(getLogger())
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
	Use:   "list-authn",
	Short: "List all authn methods on a specific application",
	Long: `List all authentication methods on a specific application
	that the user logged on can read from PAS REST API.
	
	Example Usage:
	$ cybr applications list-authn -a AppID`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfigWithLogger(getLogger())
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

var addApplicationCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an application",
	Long: `Add an application to PAS.
	
	Example Usage:
	$ cybr applications add -a AppID -l "\\"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		newApplication := requests.AddApplication{
			Application: requests.Application{
				AppID:               AppID,
				Location:            Location,
				Description:         Desc,
				AccessPermittedFrom: AccessPermittedFrom,
				AccessPermittedTo:   AccessPermittedTo,
				ExpirationDate:      ExpirationDate,
				Disabled:            Disabled,
				BusinessOwnerFName:  BusinessOwnerFName,
				BusinessOwnerLName:  BusinessOwnerLName,
				BusinessOwnerEmail:  BusinessOwnerEmail,
				BusinessOwnerPhone:  BusinessOwnerPhone,
			},
		}

		err = client.AddApplication(newApplication)
		if err != nil {
			log.Fatalf("Failed to add application '%s'. %s", newApplication.Application.AppID, err)
			return
		}

		fmt.Printf("Successfully created application '%s'\n", newApplication.Application.AppID)
	},
}

var deleteApplicationCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an application",
	Long: `Delete an application to PAS.
	
	Example Usage:
	$ cybr applications delete -a AppID`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		err = client.DeleteApplication(AppID)
		if err != nil {
			log.Fatalf("Failed to delete application '%s'. %s", AppID, err)
			return
		}

		fmt.Printf("Deleted application '%s'\n", AppID)
	},
}

var addApplicationAuthenticationMethodCmd = &cobra.Command{
	Use:   "add-authn",
	Short: "Add an authentication method to an application",
	Long: `Add an authentication method to an application to PAS.
	
	Example Usage:
	$ cybr applications add-authn -a AppID -t path -v /some/path`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		newAppAuthnMethod := requests.AddApplicationAuthentication{
			Authentication: requests.ApplicationAuthenticationMethod{
				AuthType:             AuthType,
				AuthValue:            AuthValue,
				IsFolder:             IsFolder,
				AllowInternalScripts: AllowInternalScripts,
			},
		}

		err = client.AddApplicationAuthenticationMethod(AppID, newAppAuthnMethod)
		if err != nil {
			log.Fatalf("Failed to add application authentication method to application '%s'. %s", AppID, err)
			return
		}

		fmt.Printf("Added application authentication method to '%s'\n", AppID)
	},
}

var deleteApplicationAuthenticationMethodCmd = &cobra.Command{
	Use:   "delete-authn",
	Short: "Delete an authentication method of an application",
	Long: `Delete an authentication method of an application to PAS.
	
	Example Usage:
	$ cybr applications delete-authn -a AppID -i 1`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		err = client.DeleteApplicationAuthenticationMethod(AppID, AppAuthnMethodID)
		if err != nil {
			log.Fatalf("Failed to delete application authentication method to application '%s'. %s", AppID, err)
			return
		}

		fmt.Printf("Deleted application authentication method '%s' '%s'\n", AppID, AppAuthnMethodID)
	},
}

func init() {
	// List applications
	listApplicationsCmd.Flags().StringVarP(&Location, "location", "l", "\\", "Location of the application in EPV")

	// List application methods
	listMethodsCmd.Flags().StringVarP(&AppID, "app-id", "a", "", "Application identity to filter request on")
	listMethodsCmd.MarkFlagRequired("app-id")

	// Add application
	addApplicationCmd.Flags().StringVarP(&AppID, "app-id", "a", "", "Application ID")
	addApplicationCmd.MarkFlagRequired("app-id")
	addApplicationCmd.Flags().StringVarP(&Location, "location", "l", "", "Application location")
	addApplicationCmd.MarkFlagRequired("location")
	addApplicationCmd.Flags().StringVarP(&Desc, "description", "d", "", "Application description")
	addApplicationCmd.Flags().IntVarP(&AccessPermittedFrom, "access-permitted-from", "f", 0, "Access permitted for the application. e.g. 0-23")
	addApplicationCmd.Flags().IntVarP(&AccessPermittedTo, "access-permitted-to", "t", 23, "Access permitted to the application. e.g. 0-23")
	addApplicationCmd.Flags().StringVarP(&ExpirationDate, "expiration-date", "e", "", "When application will expire")
	addApplicationCmd.Flags().StringVarP(&Disabled, "disabled", "i", "", "Disable the application. e.g. yes/no")
	addApplicationCmd.Flags().StringVarP(&BusinessOwnerFName, "business-owner-first-name", "r", "", "Application business owner first name")
	addApplicationCmd.Flags().StringVarP(&BusinessOwnerLName, "business-owner-last-name", "n", "", "Application business owner lasty name")
	addApplicationCmd.Flags().StringVarP(&BusinessOwnerEmail, "business-owner-email", "m", "", "Application business owner email")
	addApplicationCmd.Flags().StringVarP(&BusinessOwnerPhone, "business-owner-phone", "p", "", "Application business owner phone")

	// Delete application
	deleteApplicationCmd.Flags().StringVarP(&AppID, "app-id", "a", "", "Application ID")
	deleteApplicationCmd.MarkFlagRequired("app-id")

	// Add Application Authentication Method
	addApplicationAuthenticationMethodCmd.Flags().StringVarP(&AppID, "app-id", "a", "", "Application ID")
	addApplicationAuthenticationMethodCmd.MarkFlagRequired("app-id")
	addApplicationAuthenticationMethodCmd.Flags().StringVarP(&AuthType, "auth-type", "t", "", "Application authentication method type")
	addApplicationAuthenticationMethodCmd.MarkFlagRequired("auth-type")
	addApplicationAuthenticationMethodCmd.Flags().StringVarP(&AuthValue, "auth-value", "v", "", "Application authentication method value")
	addApplicationAuthenticationMethodCmd.MarkFlagRequired("auth-value")
	addApplicationAuthenticationMethodCmd.Flags().BoolVarP(&IsFolder, "is-folder", "f", false, "Application is folder")
	addApplicationAuthenticationMethodCmd.Flags().BoolVarP(&AllowInternalScripts, "allow-internal-scripts", "s", false, "Allow internal scripts")

	// Delete Application Authentication Method
	deleteApplicationAuthenticationMethodCmd.Flags().StringVarP(&AppID, "app-id", "a", "", "Application ID")
	deleteApplicationAuthenticationMethodCmd.MarkFlagRequired("app-id")
	deleteApplicationAuthenticationMethodCmd.Flags().StringVarP(&AppAuthnMethodID, "auth-method-id", "i", "", "Application authentication method ID to be deleted")
	deleteApplicationAuthenticationMethodCmd.MarkFlagRequired("auth-method-id")

	applicationsCmd.AddCommand(listApplicationsCmd)
	applicationsCmd.AddCommand(listMethodsCmd)
	applicationsCmd.AddCommand(addApplicationCmd)
	applicationsCmd.AddCommand(deleteApplicationCmd)
	applicationsCmd.AddCommand(addApplicationAuthenticationMethodCmd)
	applicationsCmd.AddCommand(deleteApplicationAuthenticationMethodCmd)

	rootCmd.AddCommand(applicationsCmd)
}
