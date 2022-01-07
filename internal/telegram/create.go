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
		msg.Text = "Internal error, please try again later..."
		b.api.Send(msg)
		return
	}

	if err == nil {
		msg.Text = "Email already created: " + email.Name
		b.api.Send(msg)
		return
	}

	if err = b.repo.CreateUserState(ctx, repository.UserState{
		ChatID: update.Message.Chat.ID,
		State:  repository.StateEmailCreationStep2,
	}); err != nil {
		log.Println("error: " + err.Error())
		msg.Text = "Internal error, please try again later..."
		b.api.Send(msg)
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
	// TODO: check if email has been already created (unique)
	// TODO: if ok -> return message ("email has been created") and create record in DB
	// TODO: if !ok -> return message ("something gone wrong")

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	email, err := b.repo.GetEmailByName(ctx, update.Message.Text)
	if err != nil && err != database.ErrNotFound {
		log.Println("error: " + err.Error())
		msg.Text = "Internal error, please try again later..."
		b.api.Send(msg)
		return
	}

	if err == nil {
		msg.Text = fmt.Sprintf("Email <%s> already exists, please choose a new one.", email.Name)
		b.api.Send(msg)
		return
	}

	if err = b.repo.CreateEmail(ctx, repository.EMail{
		ChatID: update.Message.Chat.ID,
		Name:   update.Message.Text,
	}); err != nil {
		log.Println("error: " + err.Error())
		msg.Text = "Internal error, please try again later..."
		b.api.Send(msg)
		return
	}

	if err = b.repo.DeleteUserState(ctx, repository.UserState{
		ChatID: update.Message.Chat.ID,
		State:  repository.StateEmailCreationStep2,
	}); err != nil {
		log.Println("error: " + err.Error())
		msg.Text = "Internal error, please try again later..."
		b.api.Send(msg)
		return
	}

	msg.Text = fmt.Sprintf("Email <%s> has been created succesfully.", update.Message.Text)

	if _, err = b.api.Send(msg); err != nil {
		log.Panic(err)
	}
}
