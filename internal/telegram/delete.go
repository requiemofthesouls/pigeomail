package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	"pigeomail/database"
	"pigeomail/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const deleteCommand = "delete"

func (b *Bot) handleDeleteCommandStep1(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	email, err := b.repo.GetEmailByChatID(ctx, update.Message.Chat.ID)
	if err != nil {
		if err == database.ErrNotFound {
			msg.Text = "There's no created email, use /create." + email.Name
			_, _ = b.api.Send(msg)
			return
		}

		log.Println("error: " + err.Error())
		b.internalErrorResponse(update.Message.Chat.ID)
		return
	}

	if err = b.repo.CreateUserState(ctx, repository.UserState{
		ChatID: update.Message.Chat.ID,
		State:  repository.StateDeleteEmailStep1,
	}); err != nil {
		log.Println("error: " + err.Error())
		b.internalErrorResponse(update.Message.Chat.ID)
		return
	}

	msg.Text = fmt.Sprintf("Type 'yes' if you want to delete your email: <%s>", email.Name)
	if _, err = b.api.Send(msg); err != nil {
		log.Panic(err)
	}
}

func (b *Bot) handleDeleteCommandStep2(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if update.Message.Text != "yes" {
		msg.Text = "Exiting from delete mode..."

		if err := b.repo.DeleteUserState(ctx, repository.UserState{
			ChatID: update.Message.Chat.ID,
			State:  repository.StateDeleteEmailStep1,
		}); err != nil {
			log.Println("error: " + err.Error())
			b.internalErrorResponse(update.Message.Chat.ID)
			return
		}

		_, _ = b.api.Send(msg)
		return
	}

	err := b.repo.DeleteEmail(ctx, update.Message.Chat.ID)
	if err != nil {
		log.Println("error: " + err.Error())
		b.internalErrorResponse(update.Message.Chat.ID)
		return
	}

	if err = b.repo.DeleteUserState(ctx, repository.UserState{
		ChatID: update.Message.Chat.ID,
		State:  repository.StateDeleteEmailStep1,
	}); err != nil {
		log.Println("error: " + err.Error())
		b.internalErrorResponse(update.Message.Chat.ID)
		return
	}

	msg.Text = "Email has been deleted successfully."

	if _, err = b.api.Send(msg); err != nil {
		log.Panic(err)
	}
}
