package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pigeomail/database"
	"pigeomail/internal/repository"
	"pigeomail/internal/telegram"
	"pigeomail/rabbitmq"
)

// tgBotCmd represents the tgBot command
var tgBotCmd = &cobra.Command{
	Use:   "tg_bot",
	Short: "Telegram bot which handles user input",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var tgCfg *telegram.Config
		if err = viper.UnmarshalKey("telegram", &tgCfg); err != nil {
			return err
		}

		var dbCfg *database.Config
		if err = viper.UnmarshalKey("database", &dbCfg); err != nil {
			return err
		}

		var repo repository.IEmailRepository
		if repo, err = repository.NewMongoRepository(dbCfg); err != nil {
			return err
		}

		var rmqCfg *rabbitmq.Config
		if err = viper.UnmarshalKey("rabbitmq", &rmqCfg); err != nil {
			return err
		}

		var bot *telegram.Bot
		if bot, err = telegram.NewTGBot(tgCfg, rmqCfg, repo); err != nil {
			return err
		}

		bot.Run()
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
