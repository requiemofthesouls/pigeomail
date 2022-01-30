package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"pigeomail/internal/config"
	"pigeomail/pkg/logger"
)

var configPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pigeomail",
	Short: "Service which provides securely personal temporary email addresses",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(logger.Init, initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file (default is $HOME/config.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	err := config.Init(configPath)
	cobra.CheckErr(err)
}
