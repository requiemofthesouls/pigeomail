package rabbitmq

import (
	"github.com/streadway/amqp"
)

type Config struct {
	DSN string `yaml:"dsn"`
}

type client struct {
	ch *amqp.Channel
}

func NewRMQConnection(dsn string) (*amqp.Connection, error) {
	return amqp.Dial(dsn)
}
