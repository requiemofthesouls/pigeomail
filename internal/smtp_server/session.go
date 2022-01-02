package smtp_server

import (
	"errors"
	"io"
	"log"

	"github.com/DusanKasan/parsemail"
	"github.com/emersion/go-smtp"
	"pigeomail/rabbitmq"
)

// A Session is returned after EHLO.
type Session struct {
	publisher rabbitmq.IRMQEmailPublisher
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

	err = s.publisher.PublishIncomingEmail(email)
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
