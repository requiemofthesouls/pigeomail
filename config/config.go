package config

import (
	"github.com/spf13/viper"

	"pigeomail/internal/smtp_client"
	"pigeomail/internal/smtp_server"
	"pigeomail/internal/telegram"
)

type Config struct {
	SMTP struct {
		Client smtp_client.Config `yaml:"client"`
		Server smtp_server.Config `yaml:"server"`
	} `yaml:"smtp"`
	Telegram telegram.Config `yaml:"telegram"`
}

func Get() (cfg *Config, err error) {
	err = viper.Unmarshal(&cfg)
	return cfg, err
}
