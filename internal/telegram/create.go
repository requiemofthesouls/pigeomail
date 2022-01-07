package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const createCommand = "create"

func (b *Bot) handleCreateCommand(update *tgbotapi.Update) {
	// TODO: check if email has been already created for that user
	// TODO: if ok -> return message ("enter mailbox name")
	// TODO: if !ok -> return message ("your email already created")

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = createCommand + " command in development, stay tuned..."

	if _, err := b.api.Send(msg); err != nil {
		log.Panic(err)
	}
}

func (b *Bot) handleCreateCommandCallback(update *tgbotapi.Update) {
	// TODO: check if email has been already created (unique)
	// TODO: if ok -> return message ("email has been created") and create record in DB
	// TODO: if !ok -> return message ("something gone wrong")

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = createCommand + " command in development, stay tuned..."

	if _, err := b.api.Send(msg); err != nil {
		log.Panic(err)
	}
}
