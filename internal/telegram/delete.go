package telegram

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/requiemofthesouls/pigeomail/internal/repository/entity"
)

const deleteCommand = "delete"

func (b *Bot) handleDeleteCommandStep1(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	var err error
	var email entity.EMail
	if email, err = b.repo.PrepareDeleteEmail(ctx, update.Message.Chat.ID); err != nil {
		b.handleError(err, update.Message.Chat.ID)
		return
	}

	msg.Text = fmt.Sprintf("type 'yes' if you want to delete your email: <%s>", email.Name)
	_, _ = b.api.Send(msg)
}

func (b *Bot) handleDeleteCommandStep2(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if update.Message.Text != "yes" {
		msg.Text = "exiting from delete mode..."

		if err := b.repo.CancelDeleteEmail(ctx, update.Message.Chat.ID); err != nil {
			b.handleError(err, update.Message.Chat.ID)
			return
		}

		_, _ = b.api.Send(msg)
		return
	}

	if err := b.repo.DeleteEmail(ctx, update.Message.Chat.ID); err != nil {
		b.handleError(err, update.Message.Chat.ID)
		return
	}

	msg.Text = "email has been deleted successfully"
	_, _ = b.api.Send(msg)
}
