package postgres

import (
	"fmt"
	"net/url"
)

type Config struct {
	Host               string `mapstructure:"host"`
	Port               int32  `mapstructure:"port"`
	Username           string `mapstructure:"username"`
	Password           string `mapstructure:"password"`
	Database           string `mapstructure:"database"`
	MaxConns           int32  `mapstructure:"maxConns"`
	MaxConnLifetimeSec int32  `mapstructure:"maxConnLifetimeSec"`
	MaxConnIdleTimeSec int32  `mapstructure:"maxConnIdleTimeSec"`
}

func (c *Config) setDefaultValues() {
	if c.MaxConns == 0 {
		c.MaxConns = 4
	}
	if c.MaxConnLifetimeSec == 0 {
		c.MaxConnLifetimeSec = 90
	}

	if c.MaxConnIdleTimeSec == 0 {
		c.MaxConnIdleTimeSec = 10
	}
}

func (c *Config) dsn() string {
	var (
		dsn = &url.URL{
			Scheme: "postgresql",
			Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
			Path:   c.Database,
		}
		q = dsn.Query()
	)

	q.Add("sslmode", "disable")
	dsn.RawQuery = q.Encode()

	if c.Username == "" {
		return dsn.String()
	}

	if c.Password == "" {
		dsn.User = url.User(c.Username)

		return dsn.String()
	}

	dsn.User = url.UserPassword(c.Username, c.Password)

	return dsn.String()
}
