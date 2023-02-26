package rabbitmq

import (
	"github.com/streadway/amqp"
	"go.uber.org/zap"

	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"
	rmqDef "github.com/requiemofthesouls/pigeomail/pkg/modules/rabbitmq/def"
)

type (
	Consumer interface {
		Consume(queue string, handler func(msg *amqp.Delivery)) error
	}
	consumer struct {
		rmq rmqDef.Wrapper
		l   logDef.Wrapper
	}
)

func NewConsumer(rmq rmqDef.Wrapper, l logDef.Wrapper) Consumer {
	return &consumer{rmq: rmq, l: l}
}

func (r *consumer) Consume(queue string, handler func(msg *amqp.Delivery)) (err error) {
	var l = r.l.With(zap.String("queue", queue))
	l.Info("starting RMQ consumer")

	if _, err = r.rmq.QueueDeclare(
		queue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	); err != nil {
		return err
	}

	var messages <-chan amqp.Delivery
	if messages, err = r.rmq.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	); err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			d := d
			l.Debug("RMQ Start message processing", zap.String("message_id", d.MessageId))
			handler(&d)
			l.Debug("RMQ End message processing", zap.String("message_id", d.MessageId))
		}
	}()

	l.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}
