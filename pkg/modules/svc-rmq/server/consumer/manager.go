package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/channel"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/connection"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/internal"
	"google.golang.org/protobuf/proto"
)

func NewManager(
	log logger.Wrapper,
	conn connection.Manager,
	handler QueueHandler,
	interceptors []Interceptor,
	queueHandlerSettings *QueueHandlerSettings,
	availableProducts []string,
) (*manager, error) {
	man := &manager{
		log:                       log,
		handler:                   handler,
		interceptor:               noOpInterceptor,
		exchangeName:              queueHandlerSettings.ExchangeName,
		routingKey:                queueHandlerSettings.RoutingKey,
		queueName:                 queueHandlerSettings.QueueName,
		deadQueueName:             fmt.Sprintf("%s.dead", queueHandlerSettings.QueueName),
		delayExchangeName:         fmt.Sprintf("%s.delay", queueHandlerSettings.QueueName),
		delayTTL:                  queueHandlerSettings.DelayTTL,
		delayMaxNumFailedAttempts: queueHandlerSettings.DelayNumFailedAttempts,
		qos:                       queueHandlerSettings.QOS,
		numConsumers:              queueHandlerSettings.NumConsumers,
		availableProducts:         availableProducts,
	}

	man.ctx, man.ctxCancelFunc = context.WithCancel(context.Background())

	var err error
	if man.channel, err = channel.NewManager(man.log, conn); err != nil {
		return nil, fmt.Errorf("channel create: %v", err)
	}

	if len(interceptors) > 0 {
		chainedInterceptor := func(ctx context.Context, msg *Message, handler Handler) error {
			for i := len(interceptors) - 1; i >= 0; i-- {
				currentInterceptor, currentHandler := interceptors[i], handler
				handler = func(ctx context.Context, msg *Message) error {
					return currentInterceptor(ctx, msg, currentHandler)
				}
			}
			return handler(ctx, msg)
		}
		man.interceptor = chainedInterceptor
	}

	return man, nil
}

type (
	MapQueueHandlers map[string]QueueHandlerItem

	QueueHandlerItem struct {
		ExchangeName string
		RoutingKey   string
		Handler      QueueHandler
	}

	QueueHandlerSettings struct {
		QueueName              string
		ExchangeName           string
		RoutingKey             string
		QOS                    uint8
		NumConsumers           uint8
		DelayNumFailedAttempts uint8
		DelayTTL               int32
	}

	QueueHandler interface {
		Handle(ctx context.Context, dec func(message interface{}) error) error
	}

	Interceptor func(ctx context.Context, msg *Message, handler Handler) error

	Handler func(ctx context.Context, msg *Message) error

	Message struct {
		amqp.Delivery
	}
)

type (
	Manager interface {
		Start()
		Close()
	}

	manager struct {
		log logger.Wrapper

		ctx           context.Context
		ctxCancelFunc context.CancelFunc

		channel channel.Manager

		handler     QueueHandler
		interceptor Interceptor

		exchangeName  string
		routingKey    string
		queueName     string
		deadQueueName string

		delayExchangeName         string
		delayTTL                  int32
		delayMaxNumFailedAttempts uint8

		qos uint8

		numConsumers uint8

		availableProducts []string
	}
)

func (m *manager) Start() {
	if err := m.start(); err != nil {
		m.log.Error("consumers not started", logger.Error(err))
	}

	m.channel.Close()
}

func (m *manager) Close() {
	if m.channel.IsClosed() {
		return
	}

	m.ctxCancelFunc()
}

func (m *manager) start() error {
	var err error
	if err = m.channel.Qos(channel.QosOptions{
		PrefetchCount: int(m.qos),
	}); err != nil {
		return fmt.Errorf("set qos: %v", err)
	}

	if err = m.channel.DeclareExchange(channel.ExchangeOptions{
		Name:      m.exchangeName,
		Kind:      channel.ExchangeDirect,
		IsDurable: true,
	}); err != nil {
		return fmt.Errorf("declare exchange %s: %v", m.exchangeName, err)
	}

	if err = m.channel.DeclareQueue(channel.QueueOptions{
		Name:      m.queueName,
		IsDurable: true,
	}); err != nil {
		return fmt.Errorf("declare exchange: %v", err)
	}

	if err = m.channel.DeclareQueue(channel.QueueOptions{
		Name:      m.deadQueueName,
		IsDurable: true,
	}); err != nil {
		return fmt.Errorf("declare exchange: %v", err)
	}

	if err = m.channel.BindQueueWithExchange(channel.BindingOptions{
		QueueName:    m.queueName,
		RoutingKey:   m.routingKey,
		ExchangeName: m.exchangeName,
	}); err != nil {
		return fmt.Errorf("declare queue: %v", err)
	}

	if err = m.channel.DeclareExchange(channel.ExchangeOptions{
		Name:      m.delayExchangeName,
		Kind:      "x-delayed-message",
		IsDurable: true,
		Args: channel.Table{
			"x-delayed-type": "direct",
		},
	}); err != nil {
		return fmt.Errorf("declare delay exchange %s: %v", m.delayExchangeName, err)
	}

	if err = m.channel.BindQueueWithExchange(channel.BindingOptions{
		QueueName:    m.queueName,
		ExchangeName: m.delayExchangeName,
	}); err != nil {
		return fmt.Errorf("declare queue: %v", err)
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < int(m.numConsumers); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var msgs <-chan amqp.Delivery
			if msgs, err = m.channel.Consume(channel.ConsumeOptions{
				QueueName: m.queueName,
				Name:      fmt.Sprintf("%s-%d", m.queueName, i),
			}); err != nil {
				m.log.Error("consumer not created", logger.Error(err))
				return
			}

			m.log.Debug("consumer created")
			m.handlerGoroutine(msgs)
			m.log.Debug("consumer closed")
		}(i)
	}
	wg.Wait()

	return nil
}

func (m *manager) handlerGoroutine(msgs <-chan amqp.Delivery) {
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				continue
			}

			if m.channel.IsClosed() {
				return
			}

			log := m.log.With(
				logger.String(logger.KeyRequestID, msg.MessageId),
				logger.String(logger.KeyRMQExchange, msg.Exchange),
				logger.String(logger.KeyRMQRoutingKey, msg.RoutingKey),
			)

			if err := m.interceptor(m.ctx, &Message{msg}, func(ctx context.Context, msg *Message) error {
				return m.handler.Handle(
					ctx,
					func(pMessage interface{}) error {
						switch msg.ContentType {
						case internal.ContentTypeApplicationProtobuf:
							if protoMessage, isProtoMsg := pMessage.(proto.Message); isProtoMsg {
								return proto.Unmarshal(msg.Body, protoMessage)
							}
							return fmt.Errorf("message is not proto")
						default:
							return json.Unmarshal(msg.Body, pMessage)
						}
					},
				)
			}); err != nil {
				if err = m.opFallback(&channel.MsgPublishing{
					Headers:         msg.Headers,
					ContentType:     msg.ContentType,
					ContentEncoding: msg.ContentEncoding,
					DeliveryMode:    msg.DeliveryMode,
					Priority:        msg.Priority,
					CorrelationId:   msg.CorrelationId,
					ReplyTo:         msg.ReplyTo,
					Expiration:      msg.Expiration,
					MessageId:       msg.MessageId,
					Timestamp:       msg.Timestamp,
					Type:            msg.Type,
					UserId:          msg.UserId,
					AppId:           msg.AppId,
					Body:            msg.Body,
				}, err); err != nil {
					log.Error("fallback failed", logger.Error(err))
					if err = msg.Reject(true); err != nil {
						log.Error("msg not rejected", logger.Error(err))
					}
					continue
				}

				if err = msg.Ack(false); err != nil {
					log.Error("msg not acked", logger.Error(err))
				}

				continue
			}

			if err := msg.Ack(false); err != nil {
				log.Error("msg not acked", logger.Error(err))
			}
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *manager) opFallback(msgPublishing *channel.MsgPublishing, err error) error {
	msgPublishing.Headers[msgHeaderCustomReceivedErrorMessage] = err.Error()

	var numberFailedAttempts int32

	if headerCustomNumberFailedAttemptsValue, ok := msgPublishing.Headers[msgHeaderCustomNumberFailedAttempts]; ok {
		if numberFailedAttempts, ok = headerCustomNumberFailedAttemptsValue.(int32); !ok {
			return m.channel.PublishToQueue(m.ctx, channel.PublishToQueueOptions{
				QueueName: m.deadQueueName,
				Msg:       msgPublishing,
			})
		}
	}

	if uint8(numberFailedAttempts) >= m.delayMaxNumFailedAttempts {
		return m.channel.PublishToQueue(m.ctx, channel.PublishToQueueOptions{
			QueueName: m.deadQueueName,
			Msg:       msgPublishing,
		})
	}

	msgPublishing.Headers[msgHeaderCustomNumberFailedAttempts] = numberFailedAttempts + 1
	msgPublishing.Headers[msgHeaderDelay] = m.delayTTL
	return m.channel.PublishToExchange(m.ctx, channel.PublishToExchangeOptions{
		ExchangeName: m.delayExchangeName,
		Msg:          msgPublishing,
	})
}

func noOpInterceptor(ctx context.Context, msg *Message, handler Handler) error {
	return handler(ctx, msg)
}
