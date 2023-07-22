package channel

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// QueueOptions are used to configure a queue.
// A passive queue is assumed by RabbitMQ to already exist, and attempting to connect
// to a non-existent queue will cause RabbitMQ to throw an exception.
type QueueOptions struct {
	Name         string
	IsDurable    bool
	IsAutoDelete bool
	IsExclusive  bool
	IsNoWait     bool
	Args         Table
}

func (m *manager) DeclareQueue(options QueueOptions) error {
	if _, err := m.queueDeclareSafe(
		options.Name,
		options.IsDurable,
		options.IsAutoDelete,
		options.IsExclusive,
		options.IsNoWait,
		options.Args,
	); err != nil {
		return fmt.Errorf("queue declare: %v", err)
	}
	return nil
}

// queueDeclareSafe safely wraps the (*amqp.Channel).QueueDeclare method
func (m *manager) queueDeclareSafe(
	name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table,
) (amqp.Queue, error) {
	m.RLock()
	defer m.RUnlock()

	return m.channel.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)
}
