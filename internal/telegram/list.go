package telegram

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const listCommand = "list"

func (b *Bot) handleListCommand(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = listCommand + " command in development, stay tuned..."

	email, err := b.repo.GetEmailByChatID(context.TODO(), update.Message.Chat.ID)
	if err != nil {
		msg.Text = err.Error()
		b.api.Send(msg)
		return
	}

	msg.Text = email.Name

	if _, err = b.api.Send(msg); err != nil {
		log.Panic(err)
	}

}
