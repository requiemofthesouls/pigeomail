package rabbitmq

import (
	"github.com/streadway/amqp"

	"pigeomail/pkg/logger"
)

const MessageReceivedQueueName = "h.pigeomail.MessageReceived"

func NewConnection(dsn string) (connection *amqp.Connection, err error) {
	var l = logger.GetLogger()
	l.Info("building RMQ connection")

	connection, err = amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	return connection, nil
}
