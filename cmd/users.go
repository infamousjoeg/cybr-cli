package cmd

import (
	"fmt"
	"log"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/spf13/cobra"
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

		fmt.Printf("Successfully unsuspended user'%s'\n", Username)
	},
}

func init() {
	// unsuspend
	unsuspendUserCmd.Flags().StringVarP(&Username, "username", "u", "", "The user you would like to unsuspend")
	unsuspendUserCmd.MarkFlagRequired("username")

	usersCmd.AddCommand(unsuspendUserCmd)
	rootCmd.AddCommand(usersCmd)
}
