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
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.conceal.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// Joe note: This would be good to allow multiple secret providers
}

// GetCMD returns the root cmd
func GetCMD() *cobra.Command {
	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := homedir.Dir()
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}

// 		// Search config in home directory with name ".conceal" (without extension).
// 		viper.AddConfigPath(home)
// 		viper.SetConfigName(".conceal")
// 	}

// 	viper.AutomaticEnv() // read in environment variables that match

// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// }
