package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const helpCommand = "help"

func (b *tgBot) handleHelpCommand(msg *tgbotapi.MessageConfig) {
	msg.Text = `

*Bot commands*
	/create - Create new email
	/list   - Show your email
	/delete - Delete your email
	/help   - Get help message
`
	msg.ParseMode = "markdown"
}
