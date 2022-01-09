package rabbitmq

import (
	"github.com/go-logr/logr"
	"github.com/streadway/amqp"
)

type IRMQEmailConsumer interface {
	ConsumeIncomingEmail(handler func(msg amqp.Delivery))
}

func NewRMQEmailConsumer(config *Config, log logr.Logger) (IRMQEmailConsumer, error) {
	conn, err := NewRMQConnection(config.DSN)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	r := &client{ch: ch, logger: log}
	err = r.queueDeclare()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *client) ConsumeIncomingEmail(handler func(msg amqp.Delivery)) {
	r.logger.Info("Starting RMQ consumer...")

	msgs, err := r.ch.Consume(
		MessageReceivedQueueName, // queue
		"",                       // consumer
		false,                    // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)

	if err != nil {
		r.logger.Error(err, "error starting consumer")
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			r.logger.V(10).Info(
				"RMQ Start message processing",
				"queue",
				MessageReceivedQueueName,
				"message_id",
				d.MessageId,
			)

			handler(d)

			r.logger.V(10).Info(
				"RMQ End message processing",
				"queue",
				MessageReceivedQueueName,
				"message_id",
				d.MessageId,
			)
		}
	}()

	r.logger.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
