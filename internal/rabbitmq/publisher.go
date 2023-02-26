package rabbitmq

import (
	"github.com/streadway/amqp"
	"go.uber.org/zap"

	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"
	rmqDef "github.com/requiemofthesouls/pigeomail/pkg/modules/rabbitmq/def"
)

const MessageReceivedQueueName = "h.pigeomail.MessageReceived"

type (
	Publisher interface {
		Publish(queue string, msg amqp.Publishing) error
	}
	publisher struct {
		rmq rmqDef.Wrapper
		l   logDef.Wrapper
	}
)

func NewPublisher(rmq rmqDef.Wrapper, l logDef.Wrapper) Publisher {
	return &publisher{
		rmq: rmq,
		l:   l,
	}
}

func (p *publisher) Publish(queue string, msg amqp.Publishing) (err error) {
	var l = p.l.With(zap.String("queue", queue))
	l.Debug("declaring queue")
	if _, err = p.rmq.QueueDeclare(
		queue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	); err != nil {
		return err
	}

	l.Info("publishing message")
	if err = p.rmq.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		msg,
	); err != nil {
		return err
	}
	return nil
}
