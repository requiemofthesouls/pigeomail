package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func validateIncomingMessage(msg *tgbotapi.Message) bool {
	if msg == nil { // ignore any non-Message updates
		return false
	}

	if !msg.IsCommand() { // ignore any non-command Messages
		return false
	}

	return true
}
