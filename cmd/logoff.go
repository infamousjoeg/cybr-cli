package cmd

import (
	"fmt"
	"log"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/spf13/cobra"
)

var logoffCmd = &cobra.Command{
	Use:   "logoff",
	Short: "Logoff the PAS REST API",
	Long: `Logoff the PAS REST API session established by logon.
	
	Example Usage:
	$ cybr logoff`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file written to local file system
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
		}
		// Remove the config file written to local file system
		err = client.RemoveConfig()
		if err != nil {
			log.Fatalf("Failed to remove configuration file. %s", err)
		}
		// Logoff the PAS REST API
		_ = client.Logoff()

		fmt.Println("Successfully logged off PAS.")
	},
}

func init() {
	rootCmd.AddCommand(logoffCmd)
}
