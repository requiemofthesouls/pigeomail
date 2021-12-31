package smtp_server

import (
	"github.com/emersion/go-smtp"
)

func NewBackend() (b smtp.Backend, err error) {
	return &backend{}, nil
}

// The Backend implements SMTP server methods.
type backend struct {
}

func (b *backend) NewSession(_ smtp.ConnectionState, _ string) (smtp.Session, error) {
	return &Session{}, nil
}

func (b *backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return nil, nil
}

func (b *backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return b.NewSession(*state, state.Hostname)
}
