package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"pigeomail/database"
	"pigeomail/internal/fsm"
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
			b.api.Send(msg)
			return
		}

		log.Println("error: " + err.Error())
		b.internalErrorResponse(update.Message.Chat.ID)
		return
	}

	b.usersFsmManager.SendEvent(update.Message.Chat.ID, fsm.DeleteEmail)

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

		b.usersFsmManager.SendEvent(update.Message.Chat.ID, fsm.Cancel)

		b.api.Send(msg)
		return
	}

	err := b.repo.DeleteEmail(ctx, update.Message.Chat.ID)
	if err != nil {
		log.Println("error: " + err.Error())
		b.internalErrorResponse(update.Message.Chat.ID)
		return
	}

	b.usersFsmManager.SendEvent(update.Message.Chat.ID, fsm.ConfirmDeletion)

	msg.Text = "Email has been deleted successfully."

	if _, err = b.api.Send(msg); err != nil {
		log.Panic(err)
	}

}
