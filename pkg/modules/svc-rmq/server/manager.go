package server

import (
	"sync"

	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/connection"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/server/consumer"
)

func NewManager(
	log logger.Wrapper,
	conn connection.Manager,
	listenerRegistrants []ListenerRegistrant,
	interceptors []consumer.Interceptor,
	queuesHandlersSettings map[string]*QueueHandlerSettings,
	availableProducts []string,
) (*manager, error) {
	man := &manager{
		log:                    log,
		conn:                   conn,
		interceptors:           interceptors,
		consumerGroups:         make(map[string]consumer.Manager, 0),
		queuesHandlersSettings: queuesHandlersSettings,
		availableProducts:      availableProducts,
	}

	for _, listenerRegistrant := range listenerRegistrants {
		listenerRegistrant(man)
	}

	return man, nil
}

type (
	Manager interface {
		RegisterService(handlers consumer.MapQueueHandlers)
		StartAll()
		CloseAll()
	}

	manager struct {
		log                    logger.Wrapper
		conn                   connection.Manager
		interceptors           []consumer.Interceptor
		groupID                string
		consumerGroups         map[string]consumer.Manager
		queuesHandlersSettings map[string]*QueueHandlerSettings
		availableProducts      []string
	}

	ListenerRegistrant func(srv Manager)

	QueueHandlerSettings struct {
		QOS                    uint8
		NumConsumers           uint8
		DelayNumFailedAttempts uint8
		DelayTTL               int32
	}
)

func (m *manager) RegisterService(handlers consumer.MapQueueHandlers) {
	for queueName, item := range handlers {
		if _, ok := m.consumerGroups[queueName]; ok {
			return
		}

		queueHandlerSettings := m.getQueueHandlerSettings(queueName)

		var (
			consumerMan consumer.Manager
			err         error
		)
		if consumerMan, err = consumer.NewManager(
			m.log,
			m.conn,
			item.Handler,
			m.interceptors,
			&consumer.QueueHandlerSettings{
				QueueName:              queueName,
				ExchangeName:           item.ExchangeName,
				RoutingKey:             item.RoutingKey,
				QOS:                    queueHandlerSettings.QOS,
				NumConsumers:           queueHandlerSettings.NumConsumers,
				DelayNumFailedAttempts: queueHandlerSettings.DelayNumFailedAttempts,
				DelayTTL:               queueHandlerSettings.DelayTTL,
			},
			m.availableProducts,
		); err != nil {
			m.log.Error("error create consumer", logger.Error(err))
			return
		}
		m.consumerGroups[queueName] = consumerMan
	}
}

func (m *manager) StartAll() {
	m.log.Info("Start rmq server")

	wg := &sync.WaitGroup{}

	for _, consumerGroup := range m.consumerGroups {
		wg.Add(1)
		go func(consumerGroup consumer.Manager) {
			consumerGroup.Start()
			wg.Done()
		}(consumerGroup)
	}

	wg.Wait()
	m.log.Info("Stop rmq server")
}

func (m *manager) CloseAll() {
	for _, cGroup := range m.consumerGroups {
		go cGroup.Close()
	}
}

func (m *manager) getQueueHandlerSettings(queueName string) *QueueHandlerSettings {
	if queueHandlerSettings, ok := m.queuesHandlersSettings[queueName]; ok {
		return queueHandlerSettings.setDefaultsValues()
	}

	return (&QueueHandlerSettings{}).setDefaultsValues()
}

func (s *QueueHandlerSettings) setDefaultsValues() *QueueHandlerSettings {
	if s.QOS < 1 {
		s.QOS = 10
	}

	if s.NumConsumers < 1 {
		s.NumConsumers = 5
	}

	if s.DelayNumFailedAttempts < 1 {
		s.DelayNumFailedAttempts = 5
	}

	if s.DelayTTL < 100 {
		s.DelayTTL = 2000
	}

	return s
}
