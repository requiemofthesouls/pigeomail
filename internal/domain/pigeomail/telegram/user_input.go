package telegram

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"pigeomail/internal/domain/pigeomail"
)

func (b *Bot) handleUserInput(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var state string
	var ok bool
	if state, ok = b.svc.GetUserState(ctx, update.Message.Chat.ID); !ok {
		return
	}

	switch state {
	case pigeomail.StateRequestedCreateEmail:
		b.handleCreateCommandStep2(update)
	case pigeomail.StateRequestedDeleteEmail:
		b.handleDeleteCommandStep2(update)
	}
}
