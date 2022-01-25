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

	_, _ = b.usersFsmManager.SendEvent(update.Message.Chat.ID, CreateEmail)

	msg.Text = "Enter your mailbox name:"

	if _, err = b.api.Send(msg); err != nil {
		log.Panic(err)
	}

}

func (b *Bot) validateMailboxName(email string) (bool, string) {
	if strings.Contains(email, "@") {
		return false, "please enter mailbox name without domain"
	}

	if _, err := mail.ParseAddress(email + "@" + b.domain); err != nil {
		return false, email + " is not a valid name for mailbox, please choose a new one"
	}

	return true, ""
}

func (b *Bot) handleCreateCommandStep2(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if ok, text := b.validateMailboxName(update.Message.Text); !ok {
		msg.Text = text
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

	_, _ = b.usersFsmManager.SendEvent(update.Message.Chat.ID, ChooseEmail)

	msg.Text = fmt.Sprintf("Email <%s> has been created successfully.", update.Message.Text+"@"+b.domain)

	if _, err = b.api.Send(msg); err != nil {
		log.Panic(err)
	}
}
