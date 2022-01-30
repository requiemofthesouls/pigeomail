package receiver

import (
	"github.com/emersion/go-smtp"
	"github.com/go-logr/logr"

	"pigeomail/internal/adapters/rabbitmq/publisher"
	"pigeomail/internal/domain/pigeomail"
	"pigeomail/pkg/logger"
)

func NewBackend(
	repo pigeomail.Storage,
	publisher publisher.Publisher,
) (b smtp.Backend, err error) {
	return &backend{
		publisher: publisher,
		repo:      repo,
		logger:    logger.GetLogger(),
	}, nil
}

// The Backend implements SMTP server methods.
type backend struct {
	publisher publisher.Publisher
	repo      pigeomail.Storage
	logger    *logr.Logger
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
