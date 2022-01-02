package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pigeomail/internal/telegram"
)

// tgBotCmd represents the tgBot command
var tgBotCmd = &cobra.Command{
	Use:   "tg_bot",
	Short: "Telegram bot which handles user input",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg *telegram.Config
		if err := viper.UnmarshalKey("telegram", &cfg); err != nil {
			return err
		}

		tgBot, err := telegram.NewTGBot(cfg)
		if err != nil {
			return err
		}
		tgBot.Run()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tgBotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tgBotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tgBotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
