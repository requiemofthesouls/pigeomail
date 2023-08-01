package telegram

import (
	"context"

	"github.com/requiemofthesouls/logger"
	"github.com/requiemofthesouls/pigeomail/internal/repository"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/telegram"
	"github.com/requiemofthesouls/pigeomail/pkg/state"
)

func NewBot(
	wrapper telegram.Wrapper,
	l logger.Wrapper,
	usersRep repository.TelegramUsers,
	state *state.State,
	smtpDomain string,
) *Bot {
	return &Bot{
		wrapper: wrapper,
		repositories: botRepositories{
			users: usersRep,
			state: state,
		},
		l:          l,
		smtpDomain: smtpDomain,
	}
}

type (
	Bot struct {
		wrapper      telegram.Wrapper
		repositories botRepositories
		l            logger.Wrapper
		smtpDomain   string
	}

	botRepositories struct {
		users repository.TelegramUsers
		state *state.State
	}
)

func (b *Bot) Start(ctx context.Context) {
	b.wrapper.Start(ctx, b.handleUserInput, b.handleCommand)
}
