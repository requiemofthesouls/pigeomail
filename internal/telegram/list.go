package telegram

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const listCommand = "list"

func (b *Bot) handleListCommand(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	email, err := b.repo.GetEmailByChatID(ctx, update.Message.Chat.ID)
	if err != nil {
		b.handleError(err, update.Message.Chat.ID)
		return
	}

	msg.Text = fmt.Sprintf("Your active email: <%s>", email.Name)
	_, _ = b.api.Send(msg)
}
