package cmd

import (
	"github.com/spf13/cobra"
	"pigeomail/internal/smtp_server"
)

// receiverCmd represents the receiver command
var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "SMTP server which listens incoming messages and puts them into message queue",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return smtp_server.Run()
	},
}

func init() {
	rootCmd.AddCommand(receiverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// receiverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// receiverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
