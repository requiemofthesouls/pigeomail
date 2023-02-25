package receiver

import (
	"github.com/emersion/go-smtp"
	"github.com/go-logr/logr"

	"pigeomail/internal/adapters/rabbitmq"
	"pigeomail/internal/domain/pigeomail"
	"pigeomail/pkg/logger"
)

func NewBackend(
	repo pigeomail.Storage,
	publisher rabbitmq.Publisher,
) (b smtp.Backend, err error) {
	return &backend{
		publisher: publisher,
		repo:      repo,
		logger:    logger.GetLogger(),
	}, nil
}

// The Backend implements SMTP server methods.
type backend struct {
	publisher rabbitmq.Publisher
	repo      pigeomail.Storage
	logger    *logr.Logger
}

func (b *backend) NewSession(state *smtp.Conn) (smtp.Session, error) {
	return &session{publisher: b.publisher, repo: b.repo, logger: b.logger}, nil
}
