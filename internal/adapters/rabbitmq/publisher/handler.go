package publisher

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

func NewPublisher(conn *amqp.Connection) (rabbitmq.Publisher, error) {
	var log = logger.GetLogger()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	r := &client{ch: ch, logger: log}
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *client) Publish(queue string, msg *amqp.Publishing) (err error) {
	r.logger.Info("declaring queue", "queue", queue)
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

	r.logger.Info("publishing message", "queue", queue)
	if err = r.ch.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		*msg,
	); err != nil {
		return err
	}
	return nil
}
