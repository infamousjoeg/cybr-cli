package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var logoffCmd = &cobra.Command{
	Use:   "logoff",
	Short: "Logoff the PAS REST API",
	Long: `Logoff the PAS REST API session established by logon.
	
	Example Usage:
	$ cybr logoff`,
	Run: func(cmd *cobra.Command, args []string) {
		err := client.Logoff()
		if err != nil {
			log.Fatalf("Failed to log off. %s", err)
			return
		}

		fmt.Println("Successfully logged off PAS.")
	},
}

func init() {
	rootCmd.AddCommand(logoffCmd)
}
