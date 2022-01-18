package config

import (
	"github.com/spf13/viper"

	"pigeomail/internal/receiver"
	"pigeomail/internal/telegram"
)

type Config struct {
	SMTP struct {
		Server receiver.Config `yaml:"server"`
	} `yaml:"smtp"`
	Telegram telegram.Config `yaml:"telegram"`
}

func Get() (cfg *Config, err error) {
	err = viper.Unmarshal(&cfg)
	return cfg, err
}
