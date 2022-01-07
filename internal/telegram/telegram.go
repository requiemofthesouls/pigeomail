package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type tgBot struct {
	api     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

func NewTGBot(config *Config) (*tgBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	return &tgBot{api: bot, updates: updates}, nil
}

func (b *tgBot) handleUserInput(update *tgbotapi.Update) {
	// Create a new MessageConfig. We don't have text yet,
	// so we leave it empty.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	// Extract the command from the Message.
	switch update.Message.Command() {
	case createCommand:
		b.handleCreateCommand(&msg)
	case listCommand:
		b.handleListCommand(&msg)
	case deleteCommand:
		b.handleDeleteCommand(&msg)
	case helpCommand:
		b.handleHelpCommand(&msg)
	default:
		msg.Text = "I don't know that command"
	}

	if _, err := b.api.Send(msg); err != nil {
		log.Panic(err)
	}
}

func (b *tgBot) Run() {
	for update := range b.updates {
		if !validateIncomingMessage(update.Message) {
			continue
		}

		b.handleUserInput(&update)
	}
}
