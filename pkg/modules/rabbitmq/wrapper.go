package rabbitmq

import (
	"context"

	"github.com/streadway/amqp"

	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"
)

func New(ctx context.Context, cfg Config, l logDef.Wrapper) (_ Wrapper, err error) {
	var dsn = cfg.dsn()

	l.Debug("rabbitmq: opening connection")
	var connection *amqp.Connection
	if connection, err = amqp.Dial(dsn); err != nil {
		return nil, err
	}

	l.Debug("rabbitmq: opening channel")
	var channel *amqp.Channel
	if channel, err = connection.Channel(); err != nil {
		return nil, err
	}

	return &wrapper{
		conn: connection,
		ch:   channel,
	}, nil
}

type (
	Wrapper interface {
		Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
		QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
		Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error
		CloseChannel() error
		CloseConnection() error
	}

	wrapper struct {
		conn *amqp.Connection
		ch   *amqp.Channel
	}
)

func (w *wrapper) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return w.ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (w *wrapper) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return w.ch.Publish(exchange, key, mandatory, immediate, msg)
}

func (w *wrapper) Consume(
	queue, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args amqp.Table,
) (<-chan amqp.Delivery, error) {
	return w.ch.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (w *wrapper) CloseChannel() error {
	return w.ch.Close()
}

func (w *wrapper) CloseConnection() error {
	return w.conn.Close()
}
