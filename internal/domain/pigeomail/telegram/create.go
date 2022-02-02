package telegram

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"pigeomail/internal/fsm"

	"pigeomail/internal/domain/pigeomail"
)

const createCommand = "create"

func (b *Bot) handleCreateCommandStep1(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := b.svc.PrepareCreateEmail(ctx, update.Message.Chat.ID); err != nil {
		b.handleError(err, update.Message.Chat.ID)
		return
	}

	b.usersFsmManager.SendEvent(update.Message.Chat.ID, fsm.StartCreatingEmail)

	msg.Text = "enter your mailbox name without domain"
	_, _ = b.api.Send(msg)
}

func (b *Bot) handleCreateCommandStep2(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var email = pigeomail.EMail{
		ChatID: update.Message.Chat.ID,
		Name:   update.Message.Text + "@" + b.domain,
	}

	if err := b.svc.CreateEmail(ctx, email); err != nil {
		b.handleError(err, update.Message.Chat.ID)
		return
	}

	b.usersFsmManager.SendEvent(update.Message.Chat.ID, fsm.FinishCreatingEmail)

	msg.Text = fmt.Sprintf("email <%s> has been created successfully", update.Message.Text+"@"+b.domain)
	_, _ = b.api.Send(msg)
}
