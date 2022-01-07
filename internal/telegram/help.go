package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const helpCommand = "help"

func (b *tgBot) handleHelpCommand(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = `

*Bot commands*
	/create - Create new email
	/list   - Show your email
	/delete - Delete your email
	/help   - Get help message
`
	msg.ParseMode = "markdown"

	if _, err := b.api.Send(msg); err != nil {
		log.Panic(err)
	}
}
