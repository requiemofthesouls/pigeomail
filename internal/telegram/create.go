package telegram

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/requiemofthesouls/pigeomail/internal/repository/entity"
)

const createCommand = "create"

func (b *Bot) handleCreateCommandStep1(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := b.repo.PrepareCreate(ctx, update.Message.Chat.ID); err != nil {
		b.handleError(err, update.Message.Chat.ID)
		return
	}

	msg.Text = "enter your mailbox name without domain"
	_, _ = b.api.Send(msg)
}

func (b *Bot) handleCreateCommandStep2(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := b.repo.Create(
		ctx,
		&entity.TelegramUser{
			ChatID: update.Message.Chat.ID,
			EMail:  update.Message.Text + "@" + b.domain,
		},
	); err != nil {
		b.handleError(err, update.Message.Chat.ID)
		return
	}

	msg.Text = fmt.Sprintf("email <%s> has been created successfully", update.Message.Text+"@"+b.domain)
	_, _ = b.api.Send(msg)
}
