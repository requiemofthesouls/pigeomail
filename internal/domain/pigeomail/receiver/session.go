package receiver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/jhillyerd/enmime"
	"github.com/streadway/amqp"

	"pigeomail/internal/adapters/rabbitmq/publisher"
	"pigeomail/internal/domain/pigeomail"
	"pigeomail/pkg/client/rabbitmq"
	customerrors "pigeomail/pkg/errors"
)

// A session is returned after EHLO.
type session struct {
	publisher publisher.Publisher
	repo      pigeomail.Storage
	logger    *logr.Logger
}

type email struct {
	From        string
	To          string
	Subject     string
	ContentType string
	MessageID   string
	Date        time.Time
	Body        string
	HTML        string
}

var ErrMailNotDelivered = errors.New("mail not delivered")

func (s *session) Mail(from string, opts smtp.MailOptions) error {
	s.logger.V(10).Info("mail from:", "email", from)
	return nil
}

func (s *session) Rcpt(to string) error {
	s.logger.V(10).Info("mail to:", "email", to)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := s.repo.GetEmailByName(ctx, to); err != nil {
		if err == customerrors.ErrNotFound {
			s.logger.V(10).Info("email not found, ignoring message", "email", to)
			return ErrMailNotDelivered
		}

		s.logger.Error(err, "error GetEmailByName")
		return ErrMailNotDelivered
	}

	return nil
}

func (s *session) parseMail(r io.Reader) (m *email, err error) {
	var e *enmime.Envelope
	if e, err = enmime.ReadEnvelope(r); err != nil {
		return nil, err
	}

	reg := regexp.MustCompile(`\w+@[\w.]+`)

	var toAddr string
	if toAddr = e.GetHeader("To"); toAddr != "" {
		matches := reg.FindStringSubmatch(toAddr)
		if len(matches) < 1 {
			return nil, fmt.Errorf("fail to parse destination address: %s", toAddr)
		}

		toAddr = matches[0]
	}

	m = &email{
		From:        e.GetHeader("From"),
		To:          toAddr,
		Subject:     e.GetHeader("Subject"),
		ContentType: e.GetHeader("Content-Type"),
		MessageID:   e.GetHeader("Message-Id"),
		Body:        e.Text,
		HTML:        e.HTML,
	}

	return m, nil
}

func (s *session) Data(r io.Reader) (err error) {
	var parsedEmail *email
	if parsedEmail, err = s.parseMail(r); err != nil {
		s.logger.Error(err, "error parse parsedEmail")
		return ErrMailNotDelivered
	}

	var msg = &amqp.Publishing{
		Headers: amqp.Table{
			"from":    parsedEmail.From,
			"to":      parsedEmail.To,
			"subject": parsedEmail.Subject,
			"date":    parsedEmail.Date.Unix(),
		},
		MessageId:   uuid.New().String(),
		ContentType: parsedEmail.ContentType,
		Body:        []byte(parsedEmail.Body),
	}

	if err = s.publisher.Publish(rabbitmq.MessageReceivedQueueName, msg); err != nil {
		s.logger.Error(err, "error PublishIncomingEmail")
		return ErrMailNotDelivered
	}

	return nil
}

func (s *session) Reset() {

}

func (s *session) Logout() error {
	return nil
}
