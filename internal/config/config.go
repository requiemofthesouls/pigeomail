package config

import (
	"pigeomail/logger"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Debug bool `yaml:"debug"`

	SMTP struct {
		Server struct {
			Domain            string `yaml:"domain" env-required:"true"`
			Addr              string `yaml:"addr" env-required:"true"`
			ReadTimeout       int    `yaml:"read_timeout_seconds"`
			WriteTimeout      int    `yaml:"write_timeout_seconds"`
			MaxMessageBytes   int    `yaml:"max_message_bytes"`
			MaxRecipients     int    `yaml:"max_recipients"`
			AllowInsecureAuth bool   `yaml:"allow_insecure_auth"`
		} `yaml:"server"`
	} `yaml:"smtp"`

	Database struct {
		Host     string `yaml:"host" env-required:"true"`
		Port     string `yaml:"port" env-required:"true"`
		Username string `yaml:"username" env-required:"true"`
		Password string `yaml:"password" env-required:"true"`
		DBName   string `yaml:"db_name" env-required:"true"`
	} `yaml:"database"`

	Telegram struct {
		Token   string `yaml:"token" env-required:"true"`
		Webhook struct {
			Enabled bool   `yaml:"enabled"`
			Port    string `yaml:"port"`
			Cert    string `yaml:"cert"`
			Key     string `yaml:"key"`
		} `yaml:"webhook"`
	} `yaml:"telegram"`

	Rabbit struct {
		DSN string `yaml:"dsn" env-required:"true"`
	} `yaml:"rabbitmq"`
}

var instance *Config

func Init(configPath string) (err error) {
	l := logger.GetLogger()
	l.Info("read application config")

	instance = &Config{}
	if err = cleanenv.ReadConfig(configPath, instance); err != nil {
		help, _ := cleanenv.GetDescription(instance, nil)
		l.Info(help)
		return err
	}

	return nil
}

func GetConfig() *Config {
	return instance
}
