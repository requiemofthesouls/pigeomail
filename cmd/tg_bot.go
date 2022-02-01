package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"

	storage "pigeomail/internal/adapters/db/pigeomail"
	"pigeomail/internal/adapters/rabbitmq"
	"pigeomail/internal/adapters/rabbitmq/consumer"
	"pigeomail/internal/config"
	"pigeomail/internal/domain/pigeomail"
	"pigeomail/internal/domain/pigeomail/telegram"
	"pigeomail/pkg/client/mongodb"
	rmq "pigeomail/pkg/client/rabbitmq"
	"pigeomail/pkg/logger"
	"pigeomail/pkg/monitoring"
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

		var conn *amqp.Connection
		if conn, err = rmq.NewConnection(cfg.Rabbit.DSN); err != nil {
			return err
		}

		var cons rabbitmq.Consumer
		if cons, err = consumer.NewConsumer(conn); err != nil {
			return err
		}

		var bot *telegram.Bot
		if bot, err = telegram.NewBot(
			ctx,
			cfg.Debug,
			cfg.Telegram.Webhook.Enabled,
			cfg.Telegram.Token,
			cfg.SMTP.Server.Domain,
			cfg.Telegram.Webhook.Port,
			cfg.Telegram.Webhook.Cert,
			cfg.Telegram.Webhook.Key,
			svc,
			cons,
		); err != nil {
			return err
		}

		monitoring.InitSentry(cfg.Sentry.DSN)

		bot.Run()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tgBotCmd)
}
