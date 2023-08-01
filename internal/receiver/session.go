package receiver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
	"github.com/requiemofthesouls/logger"
	pigeomail_api_pb "github.com/requiemofthesouls/pigeomail/api/pb"
	"github.com/requiemofthesouls/pigeomail/internal/repository"
	pigeomailpb "github.com/requiemofthesouls/pigeomail/pb"
	"go.uber.org/zap"
)

// A session is returned after EHLO.
type session struct {
	publisher pigeomailpb.PublisherEventsRMQClient
	repo      repository.TelegramUsers
	logger    logger.Wrapper
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

func (s *session) AuthPlain(username, password string) error {
	return nil
}

var ErrMailNotDelivered = errors.New("mail not delivered")

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	s.logger.Debug("mail from:", zap.String("email", from))
	return nil
}

func (s *session) Rcpt(to string) error {
	var l = s.logger.With(zap.String("email", to))
	l.Debug("received email")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		isExists bool
		err      error
	)
	if isExists, err = s.repo.ExistsByEMail(ctx, to); err != nil {
		l.Error("repo.ExistsByEMail error", zap.Error(err))
		return nil
	}

	if !isExists {
		l.Debug("mailbox not found, ignoring message")
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
		s.logger.Error("error parse parsedEmail", zap.Error(err))
		return ErrMailNotDelivered
	}

	if err = s.publisher.SMTPMessageV1(
		context.Background(),
		&pigeomail_api_pb.SMTPMessageEventV1{
			From:        parsedEmail.From,
			To:          parsedEmail.To,
			Subject:     parsedEmail.Subject,
			ContentType: parsedEmail.ContentType,
			Body:        parsedEmail.Body,
			Html:        parsedEmail.HTML,
			Timestamp:   parsedEmail.Date.Unix(),
		},
	); err != nil {
		s.logger.Error("s.publisher.SMTPMessageV1", zap.Error(err))
		return ErrMailNotDelivered
	}

	return nil
}

func (s *session) Reset() {

}

func (s *session) Logout() error {
	return nil
}
