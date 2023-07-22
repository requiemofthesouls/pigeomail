package connection

import (
	"fmt"
	"net/url"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type (
	Config struct {
		Host     string
		Port     int32
		Username string
		Password string
		Params   ConfigParams
	}
	ConfigParams struct {
		ConnectionName string
		Heartbeat      time.Duration
		Locale         string
	}
)

func (c *Config) setDefaultValues() {
	if c.Host == "" {
		c.Host = "localhost"
	}

	if c.Port <= 0 {
		c.Port = 5672
	}

	if c.Username == "" {
		c.Username = "guest"
	}

	if c.Password == "" {
		c.Password = "guest"
	}

	c.Params.setDefaults()
}

func (c *Config) getURL() string {
	return (&url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(c.Username, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
	}).String()
}

func (c *Config) getAMPQConfig() amqp091.Config {
	return c.Params.getAMPQConfig()
}

func (c *ConfigParams) setDefaults() {
	if c.Heartbeat <= 0 {
		c.Heartbeat = time.Second * 20
	}

	if c.Locale == "" {
		c.Locale = "en_US"
	}
}

func (c *ConfigParams) getAMPQConfig() amqp091.Config {
	cfg := amqp091.Config{
		Heartbeat:  c.Heartbeat,
		Locale:     c.Locale,
		Properties: make(map[string]interface{}, 0),
	}

	if c.ConnectionName != "" {
		cfg.Properties.SetClientConnectionName(c.ConnectionName)
	}

	return cfg
}
