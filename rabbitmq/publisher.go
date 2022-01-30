package rabbitmq

import (
	"time"

	"github.com/streadway/amqp"
)

const MessageReceivedQueueName = "h.pigeomail.MessageReceived"

type ParsedEmail struct {
	From        string
	To          string
	Subject     string
	ContentType string
	MessageID   string
	Date        time.Time
	Body        string
	HTML        string
}

type IRMQEmailPublisher interface {
	PublishIncomingEmail(msg *ParsedEmail) error
}

func NewRMQEmailPublisher(dsn string) (IRMQEmailPublisher, error) {
	conn, err := NewRMQConnection(dsn)
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

func (r *client) queueDeclare() (err error) {
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

func (r *client) PublishIncomingEmail(msg *ParsedEmail) (err error) {
	err = r.ch.Publish(
		"",                       // exchange
		MessageReceivedQueueName, // routing key
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			Headers: amqp.Table{
				"from":    msg.From,
				"to":      msg.To,
				"subject": msg.Subject,
				"date":    msg.Date.Unix(),
			},
			ContentType: msg.ContentType,
			Body:        []byte(msg.Body),
			MessageId:   msg.MessageID,
		})
	if err != nil {
		return err
	}
	return nil
}
