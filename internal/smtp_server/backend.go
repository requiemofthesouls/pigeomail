package smtp_server

import (
	"github.com/emersion/go-smtp"
	"pigeomail/internal/repository"
	"pigeomail/rabbitmq"
)

func NewBackend(p rabbitmq.IRMQEmailPublisher, r repository.IEmailRepository) (b smtp.Backend, err error) {
	return &backend{publisher: p, repo: r}, nil
}

// The Backend implements SMTP server methods.
type backend struct {
	publisher rabbitmq.IRMQEmailPublisher
	repo      repository.IEmailRepository
}

func (b *backend) NewSession(state *smtp.ConnectionState, hostname string) (smtp.Session, error) {
	return &Session{publisher: b.publisher, repo: b.repo}, nil
}

func (b *backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return nil, nil
}

func (b *backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return b.NewSession(state, state.Hostname)
}
