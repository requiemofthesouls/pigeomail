package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/streadway/amqp"
	"pigeomail/internal/repository"
	"pigeomail/rabbitmq"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	updates  tgbotapi.UpdatesChannel
	repo     repository.IEmailRepository
	consumer rabbitmq.IRMQEmailConsumer
}

func NewTGBot(config *Config, rmqCfg *rabbitmq.Config, repo repository.IEmailRepository) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	var consumer rabbitmq.IRMQEmailConsumer
	if consumer, err = rabbitmq.NewRMQEmailConsumer(rmqCfg); err != nil {
		return nil, err
	}

	return &Bot{
		api:      bot,
		updates:  updates,
		repo:     repo,
		consumer: consumer,
	}, nil
}

func (b *Bot) handleCommand(update *tgbotapi.Update) {
	// Extract the command from the Message.
	switch update.Message.Command() {
	case createCommand:
		b.handleCreateCommandStep1(update)
	case listCommand:
		b.handleListCommand(update)
	case deleteCommand:
		b.handleDeleteCommandStep1(update)
	case helpCommand:
		b.handleHelpCommand(update)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
		if _, err := b.api.Send(msg); err != nil {
			log.Panic(err)
		}
	}

}

func incomingEmailConsumer(msg amqp.Delivery) {
	log.Printf("Received a message: %s", msg.Body)
	msg.Ack(false)
}

func (b *Bot) runConsumer() {
	b.consumer.ConsumeIncomingEmail(incomingEmailConsumer)
}

func (b *Bot) Run() {
	go b.runConsumer()

	for update := range b.updates {
		if !validateIncomingMessage(update.Message) {
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(&update)
			continue
		}

		b.handleUserInput(&update)
	}
}
