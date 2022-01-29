package cmd

import (
	"github.com/emersion/go-smtp"
	"github.com/spf13/cobra"
	"pigeomail/internal/config"
	"pigeomail/internal/receiver"
	"pigeomail/internal/repository"
	"pigeomail/logger"
	"pigeomail/rabbitmq"
)

// receiverCmd represents the receiver command
var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "SMTP server which listens incoming messages and puts them into message queue",

	RunE: func(cmd *cobra.Command, args []string) (err error) {
		l := logger.GetLogger()
		l.Info("building receiver")

		var cfg = config.GetConfig()

		var repo repository.IEmailRepository
		if repo, err = repository.NewMongoRepository(
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.DBName,
			"",
		); err != nil {
			return err
		}

		var publisher rabbitmq.IRMQEmailPublisher
		if publisher, err = rabbitmq.NewRMQEmailPublisher(cfg.Rabbit.DSN); err != nil {
			return err
		}

		var backend smtp.Backend
		if backend, err = receiver.NewBackend(repo, publisher); err != nil {
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
