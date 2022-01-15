package smtp_server

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/mail"
	"regexp"
	"time"

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
	var msg *mail.Message
	if msg, err = mail.ReadMessage(r); err != nil {
		return nil, err
	}

	var body []byte
	if body, err = ioutil.ReadAll(msg.Body); err != nil {
		return nil, err
	}

	var date time.Time
	if date, err = msg.Header.Date(); err != nil {
		return nil, err
	}

	reg, _ := regexp.Compile(`[\w]+@[\w.]+`)

	var fromAddr string
	var parsedFromAddr []*mail.Address
	if parsedFromAddr, err = msg.Header.AddressList("From"); err == nil {
		fromAddr = parsedFromAddr[0].Address
	} else {
		matches := reg.FindStringSubmatch(msg.Header.Get("From"))
		if len(matches) < 1 {
			return nil, err
		}
		fromAddr = matches[0]
		err = nil
	}

	var toAddr string
	var parsedToAddr []*mail.Address
	if parsedToAddr, err = msg.Header.AddressList("To"); err == nil {
		toAddr = parsedToAddr[0].Address
	} else {
		matches := reg.FindStringSubmatch(msg.Header.Get("To"))
		if len(matches) < 1 {
			return nil, err
		}
		toAddr = matches[0]
		err = nil
	}

	m = &rabbitmq.ParsedEmail{
		From:        fromAddr,
		To:          toAddr,
		Subject:     msg.Header.Get("Subject"),
		ContentType: msg.Header.Get("Content-Type"),
		MessageID:   msg.Header.Get("Message-Id"),
		Date:        date,
		Body:        body,
	}

	return
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
