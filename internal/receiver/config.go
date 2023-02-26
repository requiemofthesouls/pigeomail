package receiver

import (
	"fmt"
)

type ServerConfig struct {
	Domain              string `mapstructure:"domain"`
	Port                uint32 `mapstructure:"port"`
	ReadTimeoutSeconds  int    `mapstructure:"read_timeout_seconds"`
	WriteTimeoutSeconds int    `mapstructure:"write_timeout_seconds"`
	MaxMessageBytes     int    `mapstructure:"max_message_bytes"`
	MaxRecipients       int    `mapstructure:"max_recipients"`
	AllowInsecureAuth   bool   `mapstructure:"allow_insecure_auth"`
}

func (c *ServerConfig) setDefaults() {
	if c.ReadTimeoutSeconds == 0 {
		c.ReadTimeoutSeconds = 10
	}
	if c.WriteTimeoutSeconds == 0 {
		c.WriteTimeoutSeconds = 10
	}
	if c.MaxMessageBytes == 0 {
		c.MaxMessageBytes = 1024
	}
	if c.MaxRecipients == 0 {
		c.MaxRecipients = 100
	}
}

func (c *ServerConfig) getAddr() string {
	return fmt.Sprintf(":%d", c.Port)
}
