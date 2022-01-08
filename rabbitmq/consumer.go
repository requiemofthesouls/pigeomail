package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

type IRMQEmailConsumer interface {
	ConsumeIncomingEmail(handler func(msg amqp.Delivery))
}

func NewRMQEmailConsumer(config *Config) (IRMQEmailConsumer, error) {
	conn, err := NewRMQConnection(config.DSN)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	r := &client{ch: ch}
	err = r.queueDeclare()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *client) ConsumeIncomingEmail(handler func(msg amqp.Delivery)) {
	log.Println("Starting RMQ consumer...")

	msgs, err := r.ch.Consume(
		MessageReceivedQueueName, // queue
		"",                       // consumer
		false,                    // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)

	if err != nil {
		log.Println("error starting consumer: ", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
