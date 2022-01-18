package smtp_server

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/go-logr/logr"
	"pigeomail/database"
	"pigeomail/internal/repository"
	"pigeomail/rabbitmq"

	"github.com/jhillyerd/enmime"
)

// A Session is returned after EHLO.
type Session struct {
	publisher rabbitmq.IRMQEmailPublisher
	repo      repository.IEmailRepository
	logger    logr.Logger
}

var ErrMailNotDelivered = errors.New("mail not delivered")

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

func parseMail(r io.Reader) (m *rabbitmq.ParsedEmail, err error) {
	var e *enmime.Envelope
	if e, err = enmime.ReadEnvelope(r); err != nil {
		return nil, err
	}

	m = &rabbitmq.ParsedEmail{
		From:        e.GetHeader("From"),
		To:          e.GetHeader("To"),
		Subject:     e.GetHeader("Subject"),
		ContentType: e.GetHeader("Content-Type"),
		MessageID:   e.GetHeader("Message-Id"),
		Body:        e.Text,
		HTML:        e.HTML,
	}

	return m, nil
}

func (s *Session) Data(r io.Reader) (err error) {
	var msg *rabbitmq.ParsedEmail
	if msg, err = parseMail(r); err != nil {
		s.logger.Error(err, "error parse email")
		return ErrMailNotDelivered
	}

	if err = s.publisher.PublishIncomingEmail(msg); err != nil {
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
