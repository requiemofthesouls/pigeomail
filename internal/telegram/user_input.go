package telegram

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/requiemofthesouls/pigeomail/internal/repository/entity"
)

func (b *Bot) handleUserInput(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var state string
	var ok bool
	if state, ok = b.repo.GetUserState(ctx, update.Message.Chat.ID); !ok {
		var msg = tgbotapi.NewMessage(update.Message.Chat.ID, "use /help to see available commands")
		if _, err := b.api.Send(msg); err != nil {
			b.logger.Error("failed to send message", zap.Error(err))
		}
		return
	}

	switch state {
	case entity.StateRequestedCreateEmail:
		b.handleCreateCommandStep2(update)
	case entity.StateRequestedDeleteEmail:
		b.handleDeleteCommandStep2(update)
	}
}
