package telegram

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
	"github.com/requiemofthesouls/pigeomail/internal/repository/entity"
)

const createCommand = "create"

func (b *Bot) handleCreateCommandStep1(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chatID := update.Message.Chat.ID

	var (
		user *entity.TelegramUser
		err  error
	)
	if user, err = b.repositories.users.GetByChatID(ctx, chatID); err != nil {
		b.handleUnexpectedError(fmt.Errorf("repositories.users.GetByChatID error: %w", err), chatID)
		return
	}

	if user.IsExist() {
		b.sendStringToUser(chatID, "user already created: "+user.EMail)
		return
	}

	b.setStateUserRequestedEmailCreation(chatID)

	b.sendStringToUser(chatID, "enter your mailbox name without domain")
}

func (b *Bot) setStateUserRequestedEmailCreation(chatID int64) {
	b.repositories.state.Add(chatID,
		fsm.NewFSM(
			entity.StateRequestedCreateEmail,
			fsm.Events{
				{
					Name: entity.StateCreateEmail,
					Src:  []string{entity.StateRequestedCreateEmail},
					Dst:  entity.StateEmailCreated,
				},
			},
			fsm.Callbacks{},
		))
}

func (b *Bot) handleCreateCommandStep2(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &entity.TelegramUser{
		ChatID: update.Message.Chat.ID,
		EMail:  update.Message.Text + "@" + b.smtpDomain,
	}
	if err := user.ValidateEMail(); err != nil {
		b.sendStringToUser(user.ChatID, err.Error())
		return
	}

	var (
		isExist bool
		err     error
	)
	if isExist, err = b.repositories.users.ExistsByEMail(ctx, user.EMail); err != nil {
		b.handleUnexpectedError(fmt.Errorf("repositories.users.ExistsByEMail error: %w", err), user.ChatID)
		return
	}

	if isExist {
		b.sendStringToUser(user.ChatID, "user <"+user.EMail+"> already exists, please choose a new one")
		return
	}

	var (
		sm *fsm.FSM
		ok bool
	)
	if sm, ok = b.repositories.state.Get(user.ChatID); !ok {
		b.sendStringToUser(user.ChatID, "create user not requested, use /create")
		return
	}

	if err = sm.Event(ctx, entity.StateCreateEmail); err != nil {
		b.handleUnexpectedError(fmt.Errorf("sm.Event error: %w", err), user.ChatID)
		return
	}

	if err = b.repositories.users.Create(ctx, user); err != nil {
		b.handleUnexpectedError(fmt.Errorf("repositories.users.Create error: %w", err), user.ChatID)
		return
	}

	b.repositories.state.Delete(user.ChatID)

	b.sendStringToUser(user.ChatID, fmt.Sprintf("email <%s> has been created successfully", update.Message.Text+"@"+b.smtpDomain))
}
