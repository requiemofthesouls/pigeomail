package consumer

import (
	"github.com/go-logr/logr"
	"github.com/streadway/amqp"

	"pigeomail/internal/adapters/rabbitmq"
	"pigeomail/pkg/logger"
)

type client struct {
	ch     *amqp.Channel
	logger *logr.Logger
}

func NewConsumer(conn *amqp.Connection) (rabbitmq.Consumer, error) {
	var log = logger.GetLogger()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	r := &client{ch: ch, logger: log}

	return r, nil
}

func (r *client) Consume(queue string, handler func(msg *amqp.Delivery)) (err error) {
	r.logger.Info("starting RMQ consumer", "queue", queue)

	if _, err = r.ch.QueueDeclare(
		queue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	); err != nil {
		return err
	}

	var messages <-chan amqp.Delivery
	if messages, err = r.ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	); err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			d := d

			r.logger.V(10).Info(
				"RMQ Start message processing",
				"queue",
				queue,
				"message_id",
				d.MessageId,
			)

			handler(&d)

			r.logger.V(10).Info(
				"RMQ End message processing",
				"queue",
				queue,
				"message_id",
				d.MessageId,
			)
		}
	}()

	r.logger.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}
