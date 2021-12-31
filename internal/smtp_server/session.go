package smtp_server

import (
	"errors"
	"io"
	"log"

	"github.com/DusanKasan/parsemail"
	"github.com/emersion/go-smtp"
	"github.com/streadway/amqp"
)

// A Session is returned after EHLO.
type Session struct {
	ch *amqp.Channel
}

func (s *Session) AuthPlain(username, password string) error {
	if username != "username" || password != "password" {
		return errors.New("Invalid username or password")
	}
	return nil
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	log.Println("Mail from:", from)
	return nil
}

func (s *Session) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	return nil
}

func (s *Session) Data(r io.Reader) (err error) {
	var email parsemail.Email
	if email, err = parsemail.Parse(r); err != nil {
		return err
	}

	log.Printf("Data: %v\n", email)

	// Put incoming message to rabbit
	q, err := s.ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	err = s.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
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

func (s *Session) Reset() {

}

func (s *Session) Logout() error {
	return nil
}
