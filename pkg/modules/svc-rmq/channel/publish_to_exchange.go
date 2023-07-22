package channel

import (
	"context"
)

type (
	PublishToExchangeOptions struct {
		ExchangeName string
		RoutingKey   string
		Msg          *MsgPublishing
	}
)

func (m *manager) PublishToExchange(ctx context.Context, options PublishToExchangeOptions) error {
	m.RLock()
	defer m.RUnlock()

	return m.channel.PublishWithContext(
		ctx,
		options.ExchangeName,
		options.RoutingKey,
		true,
		false,
		*options.Msg,
	)
}
