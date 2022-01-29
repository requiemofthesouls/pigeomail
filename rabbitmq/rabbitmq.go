package rabbitmq

import (
	"github.com/go-logr/logr"
	"github.com/streadway/amqp"
)

type Config struct {
	DSN string `yaml:"dsn"`
}

type client struct {
	ch     *amqp.Channel
	logger *logr.Logger
}

func NewRMQConnection(dsn string) (*amqp.Connection, error) {
	return amqp.Dial(dsn)
}
