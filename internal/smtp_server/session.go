package smtp_server

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/DusanKasan/parsemail"
	"github.com/emersion/go-smtp"
	"pigeomail/internal/repository"
	"pigeomail/rabbitmq"
)

// A Session is returned after EHLO.
type Session struct {
	publisher rabbitmq.IRMQEmailPublisher
	repo      repository.IEmailRepository
}

var ErrMailNotDelivered = errors.New("mail not delivered")

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := s.repo.GetEmailByName(ctx, to); err != nil {
		log.Println("error GetEmailByName:", err)
		return ErrMailNotDelivered
	}

	return nil
}

func (s *Session) Data(r io.Reader) (err error) {
	var email parsemail.Email
	if email, err = parsemail.Parse(r); err != nil {
		log.Println("error parse email:", err)
		return ErrMailNotDelivered
	}

	if err = s.publisher.PublishIncomingEmail(email); err != nil {
		log.Println("error PublishIncomingEmail:", err)
		return ErrMailNotDelivered
	}

	return nil
}

func (s *Session) Reset() {

}

func (s *Session) Logout() error {
	return nil
}
