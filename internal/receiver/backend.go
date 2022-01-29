package receiver

import (
	"pigeomail/internal/repository"
	"pigeomail/rabbitmq"

	"github.com/emersion/go-smtp"
	"github.com/go-logr/logr"
)

func NewBackend(p rabbitmq.IRMQEmailPublisher, r repository.IEmailRepository, log logr.Logger) (b smtp.Backend, err error) {
	return &backend{publisher: p, repo: r, logger: log}, nil
}

// The Backend implements SMTP server methods.
type backend struct {
	publisher rabbitmq.IRMQEmailPublisher
	repo      repository.IEmailRepository
	logger    logr.Logger
}

func (b *backend) NewSession(state *smtp.ConnectionState, hostname string) (smtp.Session, error) {
	return &Session{publisher: b.publisher, repo: b.repo, logger: b.logger}, nil
}

func (b *backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return nil, nil
}

func (b *backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return b.NewSession(state, state.Hostname)
}
