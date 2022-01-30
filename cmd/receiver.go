package cmd

import (
	"context"

	"github.com/emersion/go-smtp"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"pigeomail/internal/adapters/rabbitmq"

	storage "pigeomail/internal/adapters/db/pigeomail"
	"pigeomail/internal/adapters/rabbitmq/publisher"
	"pigeomail/internal/config"
	"pigeomail/internal/domain/pigeomail/receiver"
	"pigeomail/pkg/client/mongodb"
	rmq "pigeomail/pkg/client/rabbitmq"
	"pigeomail/pkg/logger"
)

// receiverCmd represents the receiver command
var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "SMTP server which listens incoming messages and puts them into message queue",

	RunE: func(cmd *cobra.Command, args []string) (err error) {
		l := logger.GetLogger()
		l.Info("building receiver")

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

		var rmqConn *amqp.Connection
		if rmqConn, err = rmq.NewConnection(cfg.Rabbit.DSN); err != nil {
			return err
		}

		var pub rabbitmq.Publisher
		if pub, err = publisher.NewPublisher(rmqConn); err != nil {
			return err
		}

		var backend smtp.Backend
		if backend, err = receiver.NewBackend(s, pub); err != nil {
			return err
		}

		var r *receiver.Receiver
		if r, err = receiver.NewSMTPReceiver(
			backend,
			cfg.SMTP.Server.Addr,
			cfg.SMTP.Server.Domain,
			cfg.SMTP.Server.ReadTimeout,
			cfg.SMTP.Server.WriteTimeout,
			cfg.SMTP.Server.MaxMessageBytes,
			cfg.SMTP.Server.MaxRecipients,
			cfg.SMTP.Server.AllowInsecureAuth,
		); err != nil {
			return err
		}

		return r.Run()
	},
}

func init() {
	rootCmd.AddCommand(receiverCmd)
}
