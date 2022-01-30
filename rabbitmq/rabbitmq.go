package rabbitmq

import (
	"github.com/go-logr/logr"
	"github.com/streadway/amqp"
)

type client struct {
	ch     *amqp.Channel
	logger *logr.Logger
}

func NewRMQConnection(dsn string) (*amqp.Connection, error) {
	return amqp.Dial(dsn)
}
