package cmd

import (
	"fmt"
	"log"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/prettyprint"
	"github.com/spf13/cobra"
)

var (
	// UserID is the id of a user
	UserID int
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "User actions for PAS REST API",
	Long: `All users actions that can be taken via PAS REST API.
	
	Example Usage:
	Unsuspend a User: cybr users unsuspend -u userName`,
	Aliases: []string{"user"},
}

var unsuspendUserCmd = &cobra.Command{
	Use:   "unsuspend",
	Short: "Unsuspend a specific user",
	Long: `Activates a suspended user. It does not activate an inactive user.
	
	Example Usage:
	$ cybr users unsuspend --username userName`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		err = client.UnsuspendUser(Username)
		if err != nil {
			log.Fatalf("Failed to unsuspend user '%s'. %s", Username, err)
			return
		}

		fmt.Printf("Successfully unsuspended user '%s'\n", Username)
	},
}

var listUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "List cyberark PAS users",
	Long: `Lists cyberark PAS users.
	
	Example Usage:
	$ cybr users list --search userName --filter userType`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		query := &pasapi.ListUsersQueryParams{
			Search: Search,
			Filter: Filter,
		}

		users, err := client.ListUsers(query)
		if err != nil {
			log.Fatalf("Failed to list users. %s", err)
			return
		}

		prettyprint.PrintJSON(users)
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a cyberark PAS user",
	Long: `Delete a cyberark PAS user.
	
	Example Usage:
	$ cybr users delete --id 9`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		err = client.DeleteUser(UserID)
		if err != nil {
			log.Fatalf("Failed to delete user with id '%d'. %s", UserID, err)
			return
		}

		fmt.Printf("Succesfully deleted user with id '%d'\n", UserID)
	},
}

func init() {
	// unsuspend
	unsuspendUserCmd.Flags().StringVarP(&Username, "username", "u", "", "The user you would like to unsuspend")
	unsuspendUserCmd.MarkFlagRequired("username")

	// list
	listUsersCmd.Flags().StringVarP(&Search, "search", "s", "", "Search for the username, first name or last name of a user")
	listUsersCmd.Flags().StringVarP(&Filter, "filter", "f", "", "Filter on userType or componentUser")

	// delete
	deleteUserCmd.Flags().IntVarP(&UserID, "id", "i", 0, "The ID of the user you wish to delete")
	deleteUserCmd.MarkFlagRequired("id")

	usersCmd.AddCommand(unsuspendUserCmd)
	usersCmd.AddCommand(listUsersCmd)
	usersCmd.AddCommand(deleteUserCmd)
	rootCmd.AddCommand(usersCmd)
}
