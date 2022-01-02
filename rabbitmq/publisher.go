package rabbitmq

import (
	"github.com/DusanKasan/parsemail"
	"github.com/streadway/amqp"
)

const MessageReceivedQueueName = "h.pigeomail.MessageReceived"

type IRMQEmailPublisher interface {
	PublishIncomingEmail(email parsemail.Email) error
}

func NewRMQEmailPublisher(config *Config) (IRMQEmailPublisher, error) {
	conn, err := NewRMQConnection(config.DSN)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	r := &client{ch: ch}
	err = r.registerPublisherQueues()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *client) registerPublisherQueues() (err error) {
	if _, err = r.ch.QueueDeclare(
		MessageReceivedQueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	); err != nil {
		return err
	}

	return nil
}

func (r *client) PublishIncomingEmail(email parsemail.Email) (err error) {
	err = r.ch.Publish(
		"",                       // exchange
		MessageReceivedQueueName, // routing key
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			Headers: amqp.Table{
				"from":    email.From[0].Address,
				"to":      email.To[0].Address,
				"subject": email.Subject,
			},
			ContentType: email.ContentType,
			Body:        []byte(email.TextBody),
			MessageId:   email.MessageID,
		})
	if err != nil {
		return err
	}
	return nil
}
