package cmd

import (
	"fmt"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display current version",
	Long: `Display the current version of cybr.
	
	Example Usage:
	$ cybr version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cybr v%s\n", cybr.FullVersionName)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
