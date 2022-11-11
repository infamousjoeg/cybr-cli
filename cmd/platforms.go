package cmd

import (
	"log"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/queries"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/prettyprint"
	"github.com/spf13/cobra"
)

var (
	// Active is a flag to search for platforms that are active or not
	Active bool

	// PlatformType specifies the type of platform to list
	PlatformType string

	// PlatformName specifies the name of the platform to list
	PlatformName string
)

var platformsCmd = &cobra.Command{
	Use:   "platforms",
	Short: "Platform actions for PAS REST API",
	Long: `All platform actions that can be taken via PAS REST API.
	
	Example Usage:
	List all platforms: $ cybr platforms list
	Get a Platforms details: $ cybr platforms get -i WinDomain`,
	Aliases: []string{"platform"},
}

var listPlatformsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all platforms",
	Long: `List all platforms.
	
	Example Usage:
	$ cybr platforms list`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		query := &queries.ListPlatforms{
			Active:       Active,
			PlatformType: PlatformType,
			PlatformName: PlatformName,
		}

		apps, err := client.ListPlatforms(query)
		if err != nil {
			log.Fatalf("Failed to retrieve a list of all platforms. %s", err)
			return
		}

		prettyprint.PrintJSON(apps)
	},
}

var getPlatformsCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific platform",
	Long: `Get a specific platform from PAS REST API.
	
	Example Usage:
	$ cybr platforms get -i WinDomain`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfigWithLogger(getLogger())
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		apps, err := client.GetPlatform(PlatformID)
		if err != nil {
			log.Fatalf("Failed to retrieve account '%s'. %s", PlatformID, err)
			return
		}

		prettyprint.PrintJSON(apps)
	},
}

func init() {
	// Listing platforms
	listPlatformsCmd.Flags().BoolVarP(&Active, "active", "a", false, "Filter according to whether the platform is active or not.")
	listPlatformsCmd.Flags().StringVarP(&PlatformType, "platform-type", "t", "", "Filter according to the platform type. Valid values: Group or Regular")
	listPlatformsCmd.Flags().StringVarP(&PlatformName, "platform-name", "n", "", "Filter according to the platform name. Partial matches are supported.")

	// Getting a platform
	getPlatformsCmd.Flags().StringVarP(&PlatformID, "platform-id", "i", "", "Platform ID to list from")
	getPlatformsCmd.MarkFlagRequired("platform-id")

	// Add cmd to platform cmd
	platformsCmd.AddCommand(listPlatformsCmd)
	platformsCmd.AddCommand(getPlatformsCmd)

	// Add platforms cmd to root
	rootCmd.AddCommand(platformsCmd)
}
