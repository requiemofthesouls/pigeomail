package telegram

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"time"

	"pigeomail/internal/repository"
	"pigeomail/rabbitmq"

	"github.com/go-logr/logr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/streadway/amqp"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	updates  tgbotapi.UpdatesChannel
	repo     repository.IEmailRepository
	consumer rabbitmq.IRMQEmailConsumer
	domain   string
	logger   logr.Logger
}

func NewTGBot(
	cfg *Config,
	rmqCfg *rabbitmq.Config,
	repo repository.IEmailRepository,
	domain string,
	log logr.Logger,
) (bot *Bot, err error) {
	var tgAPI *tgbotapi.BotAPI
	if tgAPI, err = tgbotapi.NewBotAPI(cfg.Token); err != nil {
		return nil, err
	}

	tgAPI.Debug = cfg.Debug

	log.Info("authorized", "account", tgAPI.Self.UserName)

	var updates tgbotapi.UpdatesChannel
	if cfg.Webhook.Enabled {
		log.Info("starting tg_bot in webhook mode", "port", cfg.Webhook.Port)

		var whCfg tgbotapi.WebhookConfig
		whCfg, err = tgbotapi.NewWebhookWithCert(
			fmt.Sprintf("https://%s:%d/%s", domain, cfg.Webhook.Port, tgAPI.Token),
			tgbotapi.FilePath(cfg.Webhook.Cert),
		)
		if err != nil {
			return nil, err
		}

		if _, err = tgAPI.Request(whCfg); err != nil {
			return nil, err
		}

		var info tgbotapi.WebhookInfo
		if info, err = tgAPI.GetWebhookInfo(); err != nil {
			return nil, err
		}

		if info.LastErrorDate != 0 {
			log.Info("telegram callback failed", "last_error", info.LastErrorMessage)
		}

		updates = tgAPI.ListenForWebhook("/" + tgAPI.Token)

		go func() {
			err = http.ListenAndServeTLS(
				fmt.Sprintf("0.0.0.0:%d", cfg.Webhook.Port),
				cfg.Webhook.Cert,
				cfg.Webhook.Key,
				nil,
			)
			if err != nil {
				log.Error(err, "error in http.ListenAndServeTLS")
			}
		}()

	} else {
		log.Info("starting tg_bot without webhook mode")

		// delete created webhook cause
		// bot won't start in that mode if webhook was created before
		deleteWHCfg := tgbotapi.DeleteWebhookConfig{
			DropPendingUpdates: false,
		}

		if _, err = tgAPI.Request(deleteWHCfg); err != nil {
			return nil, err
		}

		updateCfg := tgbotapi.NewUpdate(0)
		updateCfg.Timeout = 60

		updates = tgAPI.GetUpdatesChan(updateCfg)
	}

	var consumer rabbitmq.IRMQEmailConsumer
	if consumer, err = rabbitmq.NewRMQEmailConsumer(rmqCfg, log); err != nil {
		return nil, err
	}

	return &Bot{
		api:      tgAPI,
		updates:  updates,
		repo:     repo,
		consumer: consumer,
		domain:   domain,
		logger:   log,
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
	case helpCommand, startCommand:
		b.handleHelpCommand(update)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
		if _, err := b.api.Send(msg); err != nil {
			b.logger.Error(err, "error send message")
		}
	}
}

func (b *Bot) incomingEmailConsumer(msg *amqp.Delivery) {
	from, ok := msg.Headers["from"]
	if !ok {
		b.logger.Error(nil, "error to extract 'from' header from message")
		_ = msg.Reject(false)
	}

	to, ok := msg.Headers["to"]
	if !ok {
		b.logger.Error(nil, "error to extract 'to' header from message")
		_ = msg.Reject(false)
	}

	subject, ok := msg.Headers["subject"]
	if !ok {
		b.logger.Error(nil, "error to extract 'subject' header from message")
		_ = msg.Reject(false)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chatID, err := b.repo.GetChatIDByEmail(ctx, to.(string))
	if err != nil {
		b.logger.Error(err, "chatID not found", "email", to)
		_ = msg.Reject(false)
	}

	textTemplate := `
<b>From:</b> %s
<b>To:</b> %s
<b>Subject:</b> %s
----------------
%s
----------------
`
	if len(msg.Body) > 4096 {
		text := fmt.Sprintf(
			textTemplate,
			html.EscapeString(from.(string)),
			html.EscapeString(to.(string)),
			html.EscapeString(subject.(string)),
			html.EscapeString(string(msg.Body[:3000])),
		)

		tgMsg := tgbotapi.NewMessage(chatID, text)
		tgMsg.ParseMode = tgbotapi.ModeHTML

		if _, err = b.api.Send(tgMsg); err != nil {
			b.logger.Error(err, "error send message")
		}

		for i := 3000; i < len(msg.Body); i += 4096 {
			y := i + 4096
			if y > len(msg.Body) {
				y = len(msg.Body)
			}

			tgMsg = tgbotapi.NewMessage(chatID, html.EscapeString(string(msg.Body[i:y])))
			tgMsg.ParseMode = tgbotapi.ModeHTML

			if _, err = b.api.Send(tgMsg); err != nil {
				b.logger.Error(err, "error send message")
			}
		}

		_ = msg.Ack(false)
		return
	}

	text := fmt.Sprintf(
		textTemplate,
		html.EscapeString(from.(string)),
		html.EscapeString(to.(string)),
		html.EscapeString(subject.(string)),
		html.EscapeString(string(msg.Body)),
	)

	tgMsg := tgbotapi.NewMessage(chatID, text)
	tgMsg.ParseMode = tgbotapi.ModeHTML

	if _, err = b.api.Send(tgMsg); err != nil {
		b.logger.Error(err, "error send message")
	}

	_ = msg.Ack(false)
}

func (b *Bot) runConsumer() {
	b.consumer.ConsumeIncomingEmail(b.incomingEmailConsumer)
}

func (b *Bot) Run() {
	go b.runConsumer()

	for update := range b.updates {
		update := update
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
