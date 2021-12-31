package smtp_server

import (
	"github.com/emersion/go-smtp"
	"public_email/internal/repository"
)

func NewBackend() (b smtp.Backend, err error) {
	return &backend{}, nil
}

// The Backend implements SMTP server methods.
type backend struct {
	repo repository.IEmailRepository
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
