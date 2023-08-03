package receiver

import (
	"github.com/emersion/go-smtp"
	rmqDef "github.com/requiemofthesouls/pigeomail/cmd/rmq/def"
	sseDef "github.com/requiemofthesouls/pigeomail/internal/sse/def"

	"github.com/requiemofthesouls/logger"
	"github.com/requiemofthesouls/pigeomail/internal/repository"
)

func NewBackend(
	repo repository.TelegramUsers,
	publisher rmqDef.PublisherEventsClient,
	logger logger.Wrapper,
	sse sseDef.Server,
) (b smtp.Backend, err error) {
	return &backend{
		publisher: publisher,
		repo:      repo,
		logger:    logger,
		sse:       sse,
	}, nil
}

// The Backend implements SMTP server methods.
type (
	Backend = smtp.Backend

	backend struct {
		repo      repository.TelegramUsers
		publisher rmqDef.PublisherEventsClient
		logger    logger.Wrapper
		sse       sseDef.Server
	}
)

func (b *backend) NewSession(state *smtp.Conn) (smtp.Session, error) {
	return &session{
		publisher: b.publisher,
		repo:      b.repo,
		logger:    b.logger,
		sse:       b.sse,
	}, nil
}
