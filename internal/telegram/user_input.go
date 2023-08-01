package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/looplab/fsm"
	"go.uber.org/zap"

	"github.com/requiemofthesouls/pigeomail/internal/repository/entity"
)

func (b *Bot) handleUnexpectedError(err error, chatID int64) {
	uid := uuid.New().String()
	b.sendStringToUser(chatID, "unexpected error, contact with support and send error code: "+uid)

	b.l.Error("unexpected error, code: "+uid, zap.Error(err))
}

func (b *Bot) handleCommand(update *tgbotapi.Update) {
	switch update.Message.Command() {
	case createCommand:
		b.handleCreateCommandStep1(update)
	case listCommand:
		b.handleListCommand(update)
	case deleteCommand:
		b.handleDeleteCommandStep1(update)
	case helpCommand, startCommand:
		b.handleHelpCommand(update)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
		if _, err := b.wrapper.Send(msg); err != nil {
			b.l.Error("error send message", zap.Error(err))
		}
	}
}

func (b *Bot) handleUserInput(update *tgbotapi.Update) {
	var (
		state *fsm.FSM
		ok    bool
	)
	if state, ok = b.repositories.state.Get(update.Message.Chat.ID); !ok {
		var msg = tgbotapi.NewMessage(update.Message.Chat.ID, "use /help to see available commands")
		if _, err := b.wrapper.Send(msg); err != nil {
			b.l.Error("failed to send message", zap.Error(err))
		}
		return
	}

	switch state.Current() {
	case entity.StateRequestedCreateEmail:
		b.handleCreateCommandStep2(update)
	case entity.StateRequestedDeleteEmail:
		b.handleDeleteCommandStep2(update)
	}
}

func (b *Bot) sendStringToUser(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, _ = b.wrapper.Send(msg)
}
