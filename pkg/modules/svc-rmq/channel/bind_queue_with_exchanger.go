package channel

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type BindingOptions struct {
	QueueName    string
	RoutingKey   string
	ExchangeName string
	IsNoWait     bool
	Args         Table
}

func (m *manager) BindQueueWithExchange(options BindingOptions) error {
	if err := m.queueBindSafe(
		options.QueueName,
		options.RoutingKey,
		options.ExchangeName,
		options.IsNoWait,
		options.Args,
	); err != nil {
		return fmt.Errorf("bind queue: %v", err)
	}

	return nil
}

// queueBindSafe safely wraps the (*amqp.Channel).QueueBind method
func (m *manager) queueBindSafe(name string, key string, exchange string, noWait bool, args amqp.Table) error {
	m.RLock()
	defer m.RUnlock()

	return m.channel.QueueBind(
		name,
		key,
		exchange,
		noWait,
		args,
	)
}
