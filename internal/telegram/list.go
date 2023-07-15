package telegram

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/requiemofthesouls/pigeomail/internal/repository/entity"
)

const listCommand = "list"

func (b *Bot) handleListCommand(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		user *entity.TelegramUser
		err  error
	)
	if user, err = b.repositories.users.GetByChatID(ctx, update.Message.Chat.ID); err != nil {
		err = fmt.Errorf("repositories.users.GetByChatID error: %w", err)
		b.handleUnexpectedError(err, update.Message.Chat.ID)
		return
	}

	if !user.IsExist() {
		msg.Text = "email not found, /create a new one"
		_, _ = b.wrapper.Send(msg)
		return
	}

	msg.Text = fmt.Sprintf("Your active email: <%s>", user.EMail)
	_, _ = b.wrapper.Send(msg)
}
