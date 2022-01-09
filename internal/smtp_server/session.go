package smtp_server

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/DusanKasan/parsemail"
	"github.com/emersion/go-smtp"
	"github.com/go-logr/logr"
	"pigeomail/database"
	"pigeomail/internal/repository"
	"pigeomail/rabbitmq"
)

// A Session is returned after EHLO.
type Session struct {
	publisher rabbitmq.IRMQEmailPublisher
	repo      repository.IEmailRepository
	logger    logr.Logger
}

var ErrMailNotDelivered = errors.New("mail not delivered")

func (s *Session) AuthPlain(username, password string) error {
	if username != "username" || password != "password" {
		return errors.New("Invalid username or password")
	}
	return nil
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	s.logger.V(10).Info("mail from:", "email", from)
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.logger.V(10).Info("mail to:", "email", to)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := s.repo.GetEmailByName(ctx, to); err != nil {
		if err == database.ErrNotFound {
			s.logger.V(10).Info("email not found, ignoring message", "email", to)
			return ErrMailNotDelivered
		}

		s.logger.Error(err, "error GetEmailByName")
		return ErrMailNotDelivered
	}

	return nil
}

func (s *Session) Data(r io.Reader) (err error) {
	var email parsemail.Email
	if email, err = parsemail.Parse(r); err != nil {
		s.logger.Error(err, "error parse email")
		return ErrMailNotDelivered
	}

	if err = s.publisher.PublishIncomingEmail(email); err != nil {
		s.logger.Error(err, "error PublishIncomingEmail")
		return ErrMailNotDelivered
	}

	return nil
}

func (s *Session) Reset() {

}

func (s *Session) Logout() error {
	return nil
}
