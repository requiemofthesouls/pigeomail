package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int32  `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

func (c *Config) getClientOpts() *options.ClientOptions {

	var dsn string
	var anonymous bool
	if c.Username == "" || c.Password == "" {
		anonymous = true
		dsn = fmt.Sprintf("mongodb://%s:%d", c.Host, c.Port)
	} else {
		dsn = fmt.Sprintf("mongodb://%s:%s@%s:%d", c.Username, c.Password, c.Host, c.Port)
	}

	// TODO: разобраться
	var clientOptions = options.Client().ApplyURI(dsn)
	if !anonymous {
		clientOptions.SetAuth(options.Credential{
			AuthSource:  "",
			Username:    c.Username,
			Password:    c.Password,
			PasswordSet: true,
		})
	}

	return clientOptions
}
