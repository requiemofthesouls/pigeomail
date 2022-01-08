package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"pigeomail/database"
)

const listCommand = "list"

func (b *Bot) handleListCommand(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	email, err := b.repo.GetEmailByChatID(ctx, update.Message.Chat.ID)
	if err != nil && err == database.ErrNotFound {
		msg.Text = "Email not found, /create a new one."
		b.api.Send(msg)
		return
	}

	if err != nil && err != database.ErrNotFound {
		log.Println("error: " + err.Error())
		b.internalErrorResponse(update.Message.Chat.ID)
		return
	}

	msg.Text = fmt.Sprintf("Your active email: <%s>", email.Name)

	if _, err = b.api.Send(msg); err != nil {
		log.Panic(err)
	}

}
