package connection

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
)

var defaultBackoffPolicy = []time.Duration{
	2 * time.Second,
	5 * time.Second,
	10 * time.Second,
	15 * time.Second,
	20 * time.Second,
	25 * time.Second,
}

func NewManager(log logger.Wrapper, cfg Config) (*manager, error) {
	cfg.setDefaultValues()

	man := &manager{
		log:           log,
		cfg:           cfg,
		Mutex:         &sync.Mutex{},
		backoffPolicy: defaultBackoffPolicy,
	}

	if err := man.connect(); err != nil {
		return nil, err
	}

	man.ctx, man.ctxCancelFunc = context.WithCancel(context.Background())

	return man, nil
}

type (
	Connection = amqp.Connection

	Manager interface {
		Get() *Connection
		Close()
		IsClosed() bool
	}

	manager struct {
		ctx           context.Context
		ctxCancelFunc context.CancelFunc

		log logger.Wrapper
		cfg Config

		*sync.Mutex
		connection *amqp.Connection

		backoffPolicy []time.Duration
	}
)

func (m *manager) Get() *Connection {
	m.Lock()
	defer m.Unlock()
	return m.connection
}

func (m *manager) Close() {
	defer m.ctxCancelFunc()
	if m.IsClosed() {
		m.log.Info("RMQ connection already closed")
		return
	}

	if err := m.connection.Close(); err != nil {
		m.log.Error("amqp connection close", logger.Error(err))
	}

	m.log.Info("RMQ connection closed")
}

func (m *manager) IsClosed() bool {
	m.Lock()
	defer m.Unlock()
	return m.connection.IsClosed()
}

func (m *manager) connect() error {
	m.Lock()
	defer m.Unlock()

	var (
		conn *amqp.Connection
		err  error
	)
	if conn, err = amqp.DialConfig(m.cfg.getURL(), m.cfg.getAMPQConfig()); err != nil {
		m.log.Error("amqp dial error", logger.Error(err))
		return fmt.Errorf("amqp dial: %v", err)
	}

	m.connection = conn
	m.log.Info("RMQ connection established")

	go m.watchNotifyClose()

	return nil
}

func (m *manager) watchNotifyClose() {
	if err := <-m.connection.NotifyClose(make(chan *amqp.Error, 1)); err != nil {
		m.log.Warn("amqp connection closed with error", logger.Error(err))
		m.reconnectLoop()
		return
	}

	m.log.Debug("amqp connection closed gracefully")
}

func (m *manager) reconnectLoop() {
	iBackoffPolicy, lenBackoffPolicy := 0, len(m.backoffPolicy)-1
	timer := time.NewTimer(m.backoffPolicy[iBackoffPolicy])
	for {
		select {
		case <-timer.C:
		case <-m.ctx.Done():
			return
		}

		if err := m.connect(); err != nil {
			if iBackoffPolicy < lenBackoffPolicy {
				iBackoffPolicy++
			}
			timer.Reset(m.backoffPolicy[iBackoffPolicy])
			continue
		}

		return
	}
}
