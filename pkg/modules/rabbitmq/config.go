package rabbitmq

import (
	"fmt"
	"net/url"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int32  `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (c *Config) dsn() string {
	var dsn = &url.URL{
		Scheme: "amqp",
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
	}

	dsn.User = url.UserPassword(c.Username, c.Password)

	return dsn.String()
}
