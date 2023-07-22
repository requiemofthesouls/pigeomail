package channel

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	PublishToQueueOptions struct {
		QueueName string
		Msg       *MsgPublishing
	}

	MsgPublishing = amqp.Publishing
)

func (m *manager) PublishToQueue(ctx context.Context, options PublishToQueueOptions) error {
	m.RLock()
	defer m.RUnlock()

	return m.channel.PublishWithContext(
		ctx,
		"",
		options.QueueName,
		true,
		false,
		*options.Msg,
	)
}
