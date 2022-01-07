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
	if state, err = b.repo.GetUserStateByChatID(ctx, update.Message.Chat.ID); err != nil {
		log.Println("error: " + err.Error())
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Internal error, please try again later...")
		b.api.Send(msg)
		return
	}

	switch state.State {
	case repository.StateEmailCreationStep2:
		b.handleCreateCommandStep2(update)
		return
	}

}
