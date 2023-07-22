package channel

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/connection"
)

func NewManager(log logger.Wrapper, connection connection.Manager) (*manager, error) {
	man := &manager{
		log:        log,
		connection: connection,
		RWMutex:    &sync.RWMutex{},
	}

	if err := man.getChannel(); err != nil {
		return nil, err
	}

	man.ctx, man.ctxCancelFunc = context.WithCancel(context.Background())

	return man, nil
}

type (
	Channel = amqp.Channel
	Table   = amqp.Table

	Manager interface {
		Close()
		IsClosed() bool

		DeclareExchange(options ExchangeOptions) error
		DeclareQueue(options QueueOptions) error
		BindQueueWithExchange(options BindingOptions) error
		Consume(options ConsumeOptions) (<-chan amqp.Delivery, error)
		Qos(options QosOptions) error
		PublishToQueue(ctx context.Context, options PublishToQueueOptions) error
		PublishToExchange(ctx context.Context, options PublishToExchangeOptions) error
	}

	manager struct {
		log logger.Wrapper

		connection connection.Manager

		*sync.RWMutex

		channel *amqp.Channel

		ctx           context.Context
		ctxCancelFunc context.CancelFunc
	}
)

func (m *manager) Close() {
	defer m.ctxCancelFunc()
	if m.channel.IsClosed() {
		m.log.Debug("amqp connection already closed")
		return
	}

	if err := m.channel.Close(); err != nil {
		m.log.Error("amqp channel close", logger.Error(err))
	}

	m.log.Debug("amqp channel closed")
}

func (m *manager) IsClosed() bool {
	m.Lock()
	defer m.Unlock()
	return m.channel.IsClosed()
}

func (m *manager) getChannel() error {
	m.Lock()
	defer m.Unlock()

	if m.connection.IsClosed() {
		m.log.Warn("ampq get channel: connection closed")
		return errors.New("get channel: connection closed")
	}

	var err error
	if m.channel, err = m.connection.Get().Channel(); err != nil {
		m.log.Error("ampq get channel", logger.Error(err))
		return fmt.Errorf("create channel: %v", err)
	}

	m.log.Debug("amqp channel created")

	go m.watchNotifyCancelOrClose()

	return nil
}

func (m *manager) watchNotifyCancelOrClose() {
	select {
	case err := <-m.channel.NotifyClose(make(chan *amqp.Error, 1)):
		if err != nil {
			m.log.Error("amqp channel closed with error", logger.Error(err))
			m.reconnectLoop()
			return
		}

		m.log.Debug("amqp channel closed gracefully")
	case cancelStr := <-m.channel.NotifyCancel(make(chan string, 1)):
		m.log.Warn(fmt.Sprintf("amqp channel cancel: %s", cancelStr))
		m.reconnectLoop()
	}
}

func (m *manager) reconnectLoop() {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
		case <-m.ctx.Done():
			return
		}

		if err := m.getChannel(); err != nil {
			continue
		}

		return
	}
}
