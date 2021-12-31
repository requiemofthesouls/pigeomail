package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// senderCmd represents the sender command
var senderCmd = &cobra.Command{
	Use:   "sender",
	Short: "SMTP client which listens messages in queue and sends them",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sender called")
	},
}

func init() {
	rootCmd.AddCommand(senderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// senderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// senderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
