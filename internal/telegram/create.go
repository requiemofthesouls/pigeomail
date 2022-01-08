package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"pigeomail/database"
	"pigeomail/internal/repository"
)

const createCommand = "create"

func (b *Bot) handleCreateCommandStep1(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	email, err := b.repo.GetEmailByChatID(ctx, update.Message.Chat.ID)
	if err != nil && err != database.ErrNotFound {
		log.Println("error: " + err.Error())
		b.internalErrorResponse(update.Message.Chat.ID)
		return
	}

	if err == nil {
		msg.Text = "Email already created: " + email.Name + "@" + b.domain
		b.api.Send(msg)
		return
	}

	if err = b.repo.CreateUserState(ctx, repository.UserState{
		ChatID: update.Message.Chat.ID,
		State:  repository.StateCreateEmailStep1,
	}); err != nil {
		log.Println("error: " + err.Error())
		b.internalErrorResponse(update.Message.Chat.ID)
		return
	}

	msg.Text = "Enter your mailbox name:"

	if _, err = b.api.Send(msg); err != nil {
		log.Panic(err)
	}

}

func validateMailboxName() bool {
	// TODO: email validation
	return true
}

func (b *Bot) handleCreateCommandStep2(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	email, err := b.repo.GetEmailByName(ctx, update.Message.Text)
	if err != nil && err != database.ErrNotFound {
		log.Println("error: " + err.Error())
		b.internalErrorResponse(update.Message.Chat.ID)
		return
	}

	if err == nil {
		msg.Text = fmt.Sprintf("Email <%s> already exists, please choose a new one.", email.Name+"@"+b.domain)
		b.api.Send(msg)
		return
	}

	if err = b.repo.CreateEmail(ctx, repository.EMail{
		ChatID: update.Message.Chat.ID,
		Name:   update.Message.Text + "@" + b.domain,
	}); err != nil {
		log.Println("error: " + err.Error())
		b.internalErrorResponse(update.Message.Chat.ID)
		return
	}

	if err = b.repo.DeleteUserState(ctx, repository.UserState{
		ChatID: update.Message.Chat.ID,
		State:  repository.StateCreateEmailStep1,
	}); err != nil {
		log.Println("error: " + err.Error())
		b.internalErrorResponse(update.Message.Chat.ID)
		return
	}

	msg.Text = fmt.Sprintf("Email <%s> has been created successfully.", update.Message.Text+"@"+b.domain)

	if _, err = b.api.Send(msg); err != nil {
		log.Panic(err)
	}
}
