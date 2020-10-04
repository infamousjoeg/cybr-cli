package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	pasapi "github.com/infamousjoeg/pas-api-go/pkg/cybr/api"
	"github.com/kataras/tablewriter"
	"github.com/landoop/tableprinter"
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
		}
		// List All Safes
		safes, err := client.ListSafes()
		if err != nil {
			log.Fatalf("Failed to retrieve a list of all safes. %s", err)
			return
		}

		// This is where I start to attempt a pretty print using tableprinter
		// https://github.com/lensesio/tableprinter#examples
		// I start by taking the returned safes whatever it's called and marshal to []byte
		bSafes, err := json.Marshal(safes)
		if err != nil {
			log.Fatalf("Failed to marshal the list of all safes. %s", err)
			return
		}

		// This is for debugging so I can see what's returned
		fmt.Printf("%s\n\n", string(bSafes))

		// This initializes the tableprinter to print to Stdout
		printer := tableprinter.New(os.Stdout)

		// Optionally, customize the table, import of the underline 'tablewriter' package is required for that.
		printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
		printer.CenterSeparator = "│"
		printer.ColumnSeparator = "│"
		printer.RowSeparator = "─"
		printer.HeaderBgColor = tablewriter.BgBlackColor
		printer.HeaderFgColor = tablewriter.FgGreenColor

		// Here is where I print the []byte type var bSafes
		// You can also do printer.PrintJSON if it's a string of legit JSON, but this seems to return
		// maps for each JSON entry in the Safes array...
		printer.Print(bSafes)

	},
}

func init() {
	safesCmd.Flags().StringVarP(&Safe, "safe", "s", "", "Safe name to filter request on")
	rootCmd.AddCommand(safesCmd)
}
