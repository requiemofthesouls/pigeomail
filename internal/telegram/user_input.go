package telegram

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"pigeomail/internal/repository"
)

func (b *Bot) handleUserInput(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var state repository.UserState
	var err error
	if state, err = b.repo.GetUserState(ctx, update.Message.Chat.ID); err != nil {
		log.Println("error: " + err.Error())
		b.internalErrorResponse(update.Message.Chat.ID)
		return
	}

	switch state.State {
	case repository.StateCreateEmailStep1:
		b.handleCreateCommandStep2(update)
	case repository.StateDeleteEmailStep1:
		b.handleDeleteCommandStep2(update)
	}

}
