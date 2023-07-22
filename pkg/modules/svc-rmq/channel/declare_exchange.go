package channel

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ExchangeOptions struct {
	Name         string
	Kind         string // possible values: empty string for default exchange or direct, topic, fanout
	IsDurable    bool
	IsAutoDelete bool
	IsInternal   bool
	IsNoWait     bool
	Args         Table
}

func (m *manager) DeclareExchange(options ExchangeOptions) error {
	if err := m.exchangeDeclareSafe(
		options.Name,
		options.Kind,
		options.IsDurable,
		options.IsAutoDelete,
		options.IsInternal,
		options.IsNoWait,
		options.Args,
	); err != nil {
		return fmt.Errorf("exchange declare: %v", err)
	}
	return nil
}

// exchangeDeclareSafe safely wraps the (*amqp.Channel).ExchangeDeclare method
func (m *manager) exchangeDeclareSafe(
	name string,
	kind string,
	durable bool,
	autoDelete bool,
	internal bool,
	noWait bool,
	args amqp.Table,
) error {
	m.RLock()
	defer m.RUnlock()

	return m.channel.ExchangeDeclare(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		args,
	)
}
