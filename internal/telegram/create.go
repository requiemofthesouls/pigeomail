package telegram

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"strings"
	"time"

	"pigeomail/database"
	"pigeomail/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func validateMailboxName(email string) bool {
	// Adding a simple stub string, as validation requires full email address
	// since user will give us just the name of inbox
	domainStub := "@pigeomail.com"
	_, err := mail.ParseAddress(email + domainStub)
	return err == nil
}

func (b *Bot) handleCreateCommandStep2(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if strings.Contains(update.Message.Text, "@") {
		msg.Text = fmt.Sprintf("please don't provide domain name, it is a mailbox, we provide domain for you. <%s> ", update.Message.Text)
		b.api.Send(msg)
		return
	}

	if !validateMailboxName(update.Message.Text) {
		msg.Text = fmt.Sprintf("<%s> is not a valid name for email inbox, please choose a new one.", update.Message.Text)
		b.api.Send(msg)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
