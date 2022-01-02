package smtp_server

import (
	"github.com/emersion/go-smtp"
	"pigeomail/rabbitmq"
)

func NewBackend(p rabbitmq.IRMQEmailPublisher) (b smtp.Backend, err error) {
	return &backend{publisher: p}, nil
}

// The Backend implements SMTP server methods.
type backend struct {
	publisher rabbitmq.IRMQEmailPublisher
}

func (b *backend) NewSession(state *smtp.ConnectionState, hostname string) (smtp.Session, error) {
	return &Session{publisher: b.publisher}, nil
}

func (b *backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return nil, nil
}

func (b *backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return b.NewSession(state, state.Hostname)
}
