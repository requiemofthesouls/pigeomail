package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const deleteCommand = "delete"

func (b *tgBot) handleDeleteCommand(msg *tgbotapi.MessageConfig) {
	msg.Text = deleteCommand + " command in development, stay tuned..."

}
