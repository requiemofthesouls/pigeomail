package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/channel"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/connection"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/internal"
)

const defaultChannel = "default"

type (
	Manager interface {
		Close()
		PublishToExchange(ctx context.Context, exchange, routingKey string, ev interface{}) error
	}

	manager struct {
		log      logger.Wrapper
		conn     connection.Manager
		channels map[string]channel.Manager

		declaredExchanges map[string]struct{}
	}
)

func NewManager(log logger.Wrapper, conn connection.Manager, routes []string) (*manager, error) {
	man := manager{
		log:               log,
		conn:              conn,
		channels:          make(map[string]channel.Manager, 0),
		declaredExchanges: make(map[string]struct{}, 0),
	}

	routes = append(routes, defaultChannel)

	for _, path := range routes {
		if _, ok := man.channels[path]; ok {
			continue
		}

		var err error
		if man.channels[path], err = channel.NewManager(log, conn); err != nil {
			return nil, fmt.Errorf("error create channel manager: %v", err)
		}
	}

	return &man, nil
}

func (m *manager) Close() {
	wg := &sync.WaitGroup{}

	for _, ch := range m.channels {
		wg.Add(1)
		go func(ch channel.Manager) {
			ch.Close()
			wg.Done()
		}(ch)
	}

	wg.Wait()
	m.log.Info("close rmq client")
}

func (m *manager) getChannel(exchange, routingKey string) channel.Manager {
	if ch, ok := m.channels[strings.Join([]string{exchange, routingKey}, "/")]; ok {
		return ch
	}

	return m.channels[defaultChannel]
}

func (m *manager) PublishToExchange(ctx context.Context, exchangeName, routingKey string, ev interface{}) error {
	ch := m.getChannel(exchangeName, routingKey)

	if _, ok := m.declaredExchanges[exchangeName]; !ok {
		if err := ch.DeclareExchange(channel.ExchangeOptions{
			Name:      exchangeName,
			Kind:      channel.ExchangeDirect,
			IsDurable: true,
		}); err != nil {
			return fmt.Errorf("exchange %s not declared: %v", exchangeName, err)
		}
		m.declaredExchanges[exchangeName] = struct{}{}
	}

	return ch.PublishToExchange(ctx, channel.PublishToExchangeOptions{
		ExchangeName: exchangeName,
		RoutingKey:   routingKey,
		Msg:          getMsgPublishing(ctx, ev),
	})
}

func getMsgPublishing(ctx context.Context, ev interface{}) *channel.MsgPublishing {
	msgBody, _ := json.Marshal(ev)

	return &channel.MsgPublishing{
		Headers: channel.Table{
			internal.HeaderRequestID: "requestID",
		},
		ContentType:  internal.ContentTypeApplicationJSON,
		DeliveryMode: channel.DeliveryModePersistent,
		MessageId:    "requestID",
		Body:         msgBody,
	}
}
