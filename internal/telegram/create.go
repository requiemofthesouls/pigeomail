package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const createCommand = "create"

func (b *tgBot) handleCreateCommand(msg *tgbotapi.MessageConfig) {
	msg.Text = createCommand + " command in development, stay tuned..."
}
