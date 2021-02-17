package cmd

import (
	"fmt"
	"os"

	"github.com/infamousjoeg/cybr-cli/pkg/logger"
	"github.com/spf13/cobra"
)

var cfgFile string

var (
	// Verbose logging
	Verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cybr",
	Short: "cybr is CyberArk's PAS command-line interface utility",
	Long: `cybr is a command-line interface utility created by CyberArk that
wraps the PAS REST API and eases the user experience for automators
and automation to easily interact with CyberArk Privileged Access
Security.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getLogger() logger.CMD {
	return logger.CMD{
		LoggerEnabled:    Verbose,
		LogHeaderEnabled: true,
		LogBodyEnabled:   true,
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "To enable verbose logging")
}

// GetCMD returns the root cmd
func GetCMD() *cobra.Command {
	return rootCmd
}
