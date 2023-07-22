package channel

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumeOptions struct {
	QueueName   string
	Name        string
	IsAutoAck   bool
	IsExclusive bool
	IsNoLocal   bool
	IsNoWait    bool
	Args        Table
}

func (m *manager) Consume(options ConsumeOptions) (<-chan amqp.Delivery, error) {
	m.RLock()
	defer m.RUnlock()

	return m.channel.Consume(
		options.QueueName,
		options.Name,
		options.IsAutoAck,
		options.IsExclusive,
		options.IsNoLocal,
		options.IsNoWait,
		options.Args,
	)
}

type QosOptions struct {
	PrefetchCount int
	PrefetchSize  int
	IsGlobal      bool
}

func (m *manager) Qos(options QosOptions) error {
	m.RLock()
	defer m.RUnlock()

	return m.channel.Qos(options.PrefetchCount, options.PrefetchSize, options.IsGlobal)
}
