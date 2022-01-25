package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleUserInput(update *tgbotapi.Update) {
	var state = b.usersFsmManager.GetState(update.Message.Chat.ID)

	switch state {
	case ChoosingEmail:
		b.handleCreateCommandStep2(update)
	case DeletingEmail:
		b.handleDeleteCommandStep2(update)
	}

}
