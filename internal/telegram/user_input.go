package telegram

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"pigeomail/pkg/logger"

	"pigeomail/internal/domain/pigeomail"
)

func (b *Bot) handleUserInput(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var l = logger.GetLogger()

	var state pigeomail.UserState
	var err error
	if state, err = b.svc.GetUserState(ctx, update.Message.Chat.ID); err != nil {
		l.Error(err, "error GetUserState")
		return
	}

	switch state.State {
	case pigeomail.StateCreateEmailStep1:
		b.handleCreateCommandStep2(update)
	case pigeomail.StateDeleteEmailStep1:
		b.handleDeleteCommandStep2(update)
	}
}
