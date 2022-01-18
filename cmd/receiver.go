package cmd

import (
	"pigeomail/database"
	"pigeomail/internal/receiver"
	"pigeomail/internal/repository"
	"pigeomail/logger"
	"pigeomail/rabbitmq"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// receiverCmd represents the receiver command
var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "SMTP server which listens incoming messages and puts them into message queue",

	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var rmqCfg *rabbitmq.Config
		if err = viper.UnmarshalKey("rabbitmq", &rmqCfg); err != nil {
			return err
		}

		var smtpCfg *receiver.Config
		if err = viper.UnmarshalKey("smtp.server", &smtpCfg); err != nil {
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

		var log = logger.New()

		var r *receiver.Receiver
		if r, err = receiver.NewSMTPReceiver(rmqCfg, smtpCfg, repo, log); err != nil {
			return err
		}

		return r.Run()
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
