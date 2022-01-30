package rabbitmq

import (
	"github.com/streadway/amqp"
)

const MessageReceivedQueueName = "h.pigeomail.MessageReceived"

type Consumer interface {
	Consume(queue string, handler func(msg *amqp.Delivery)) error
}

type Publisher interface {
	Publish(queue string, msg *amqp.Publishing) error
}
