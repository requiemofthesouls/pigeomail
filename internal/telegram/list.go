package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const listCommand = "list"

func (b *tgBot) handleListCommand(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = listCommand + " command in development, stay tuned..."

	if _, err := b.api.Send(msg); err != nil {
		log.Panic(err)
	}

}
