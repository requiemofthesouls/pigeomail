package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) internalErrorResponse(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Internal error, please try again later...")
	b.api.Send(msg)
	return
}
