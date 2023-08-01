package telegram

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
	"github.com/requiemofthesouls/pigeomail/internal/repository/entity"
)

const deleteCommand = "delete"

func (b *Bot) handleDeleteCommandStep1(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		user *entity.TelegramUser
		err  error
	)

	if user, err = b.repositories.users.GetByChatID(ctx, update.Message.Chat.ID); err != nil {
		b.handleUnexpectedError(fmt.Errorf("repositories.users.GetByChatID error: %w", err), update.Message.Chat.ID)
		return
	}

	if !user.IsExist() {
		b.sendStringToUser(update.Message.Chat.ID, "there's no created emails, use /create")
		return
	}

	b.repositories.state.Add(user.ChatID,
		fsm.NewFSM(
			entity.StateRequestedDeleteEmail,
			fsm.Events{
				{
					Name: entity.StateDeleteEmail,
					Src:  []string{entity.StateRequestedDeleteEmail},
					Dst:  entity.StateEmailDeleted,
				},
			},
			fsm.Callbacks{},
		))

	b.sendStringToUser(update.Message.Chat.ID, fmt.Sprintf("type 'yes' if you want to delete your email: <%s>", user.EMail))
}

func (b *Bot) handleDeleteCommandStep2(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if update.Message.Text != "yes" {
		var ok bool
		if _, ok = b.repositories.state.Get(update.Message.Chat.ID); !ok {
			b.sendStringToUser(update.Message.Chat.ID, "delete email not requested, use /delete")
			return
		}

		b.repositories.state.Delete(update.Message.Chat.ID)
		b.sendStringToUser(update.Message.Chat.ID, "exiting from delete mode...")
		return
	}

	if err := b.repositories.users.DeleteByChatID(ctx, update.Message.Chat.ID); err != nil {
		b.handleUnexpectedError(err, update.Message.Chat.ID)
		return
	}

	b.repositories.state.Delete(update.Message.Chat.ID)
	b.sendStringToUser(update.Message.Chat.ID, "email has been deleted successfully")
}
