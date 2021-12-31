package smtp_server

import (
	"github.com/emersion/go-smtp"
	"github.com/streadway/amqp"
)

func NewBackend(ch *amqp.Channel) (b smtp.Backend, err error) {
	return &backend{ch}, nil
}

// The Backend implements SMTP server methods.
type backend struct {
	ch *amqp.Channel
}

func (b *backend) NewSession(state *smtp.ConnectionState, hostname string) (smtp.Session, error) {
	return &Session{b.ch}, nil
}

func (b *backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return nil, nil
}

func (b *backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return b.NewSession(state, state.Hostname)
}
