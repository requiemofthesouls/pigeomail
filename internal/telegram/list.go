package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const listCommand = "list"

func (b *tgBot) handleListCommand(msg *tgbotapi.MessageConfig) {
	msg.Text = listCommand + " command in development, stay tuned..."

}
