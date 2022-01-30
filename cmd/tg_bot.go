package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"

	storage "pigeomail/internal/adapters/db/pigeomail"
	"pigeomail/internal/config"
	"pigeomail/internal/domain/pigeomail"
	"pigeomail/internal/domain/pigeomail/telegram"
	"pigeomail/pkg/client/mongodb"
	"pigeomail/pkg/logger"
)

// tgBotCmd represents the tgBot command
var tgBotCmd = &cobra.Command{
	Use:   "tg_bot",
	Short: "Telegram bot which handles user input",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var l = logger.GetLogger()
		l.Info("building tg_bot")

		var cfg = config.GetConfig()

		ctx := context.Background()

		var db *mongo.Database
		if db, err = mongodb.NewClient(
			ctx,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.DBName,
			"",
		); err != nil {
			return err
		}

		var s = storage.NewStorage(db)
		var svc = pigeomail.NewService(s)

		var bot *telegram.Bot
		if bot, err = telegram.NewTGBot(
			ctx,
			cfg,
			svc,
		); err != nil {
			return err
		}

		bot.Run()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tgBotCmd)
}
