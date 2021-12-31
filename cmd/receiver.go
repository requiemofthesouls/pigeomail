package cmd

import (
	"time"

	"github.com/emersion/go-smtp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pigeomail/internal/smtp_server"
)

// receiverCmd represents the receiver command
var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "SMTP server which listens incoming messages and puts them into message queue",

	RunE: func(cmd *cobra.Command, args []string) error {
		var s *smtp.Server
		var err error
		if s, err = build(); err != nil {
			return err
		}

		return smtp_server.Run(s)
	},
}

// build Builds smtp server with options in config
func build() (s *smtp.Server, err error) {
	var b smtp.Backend
	if b, err = smtp_server.NewBackend(); err != nil {
		return nil, err
	}

	s = smtp.NewServer(b)

	var cfg *smtp_server.Config
	if err = viper.UnmarshalKey("smtp.server", &cfg); err != nil {
		return nil, err
	}

	s.Addr = cfg.Addr
	s.Domain = cfg.Domain
	s.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Second
	s.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Second
	s.MaxMessageBytes = cfg.MaxMessageBytes * 1024
	s.MaxRecipients = cfg.MaxRecipients
	s.AllowInsecureAuth = cfg.AllowInsecureAuth

	return s, nil
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
