package receiver

import (
	"github.com/emersion/go-smtp"

	"github.com/requiemofthesouls/pigeomail/internal/repository"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
)

type Backend = smtp.Backend

func NewBackend(
	repo repository.TelegramUsers,
	//publisher rabbitmq.Publisher,
	logger logger.Wrapper,
) (b smtp.Backend, err error) {
	return &backend{
		//publisher: publisher,
		repo:   repo,
		logger: logger,
	}, nil
}

// The Backend implements SMTP server methods.
type backend struct {
	repo repository.TelegramUsers
	//publisher rabbitmq.Publisher
	logger logger.Wrapper
}

func (b *backend) NewSession(state *smtp.Conn) (smtp.Session, error) {
	return &session{
		//publisher: b.publisher,
		repo:   b.repo,
		logger: b.logger,
	}, nil
}
