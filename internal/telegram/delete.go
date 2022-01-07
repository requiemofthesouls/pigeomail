package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const deleteCommand = "delete"

func (b *Bot) handleDeleteCommand(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = deleteCommand + " command in development, stay tuned..."

	if _, err := b.api.Send(msg); err != nil {
		log.Panic(err)
	}

}
