package telegram

import (
	"context"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"pigeomail/internal/domain/pigeomail"
	customerrors "pigeomail/internal/errors"
)

func (b *Bot) handleUserInput(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var state pigeomail.UserState
	var err error
	if state, err = b.svc.GetUserState(ctx, update.Message.Chat.ID); err != nil {
		sentry.CaptureMessage("test123")
		if err != customerrors.ErrNotFound {
			log.Println("error: " + err.Error())
			sentry.CaptureException(err)
			b.handleError(err, update.Message.Chat.ID)
		}
		return
	}

	switch state.State {
	case pigeomail.StateCreateEmailStep1:
		b.handleCreateCommandStep2(update)
	case pigeomail.StateDeleteEmailStep1:
		b.handleDeleteCommandStep2(update)
	}
}
