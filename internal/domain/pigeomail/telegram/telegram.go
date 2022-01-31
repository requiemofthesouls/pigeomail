package telegram

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/http"
	"os"
	"time"

	"github.com/go-logr/logr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/streadway/amqp"

	"pigeomail/internal/adapters/rabbitmq"
	"pigeomail/internal/domain/pigeomail"
	customerrors "pigeomail/internal/errors"
	"pigeomail/pkg/logger"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	updates  tgbotapi.UpdatesChannel
	svc      pigeomail.Service
	consumer rabbitmq.Consumer
	domain   string
	logger   *logr.Logger
}

var aaAAAA = ""

func getWebhookUpdatesChan(
	tgAPI *tgbotapi.BotAPI,
	domain, port, cert, key string,
) (updates tgbotapi.UpdatesChannel, err error) {
	var log = logger.GetLogger()
	log.Info("starting tg_bot in webhook mode", "port", port)

	var whCfg tgbotapi.WebhookConfig
	whCfg, err = tgbotapi.NewWebhookWithCert(
		fmt.Sprintf("https://%s:%s/%s", domain, port, tgAPI.Token),
		tgbotapi.FilePath(cert),
	)
	if err != nil {
		log.Error(err, "fail to initialize tgbotapi.NewWebhookWithCert")
		return nil, err
	}

	if _, err = tgAPI.Request(whCfg); err != nil {
		log.Error(err, "fail set webhook")
		return nil, err
	}

	var info tgbotapi.WebhookInfo
	if info, err = tgAPI.GetWebhookInfo(); err != nil {
		log.Error(err, "fail GetWebhookInfo")
		return nil, err
	}

	if info.LastErrorDate != 0 {
		log.Info("GetWebhookInfo", "last_error", info.LastErrorMessage)
	}

	updates = tgAPI.ListenForWebhook("/" + tgAPI.Token)

	go func() {
		err = http.ListenAndServeTLS(
			fmt.Sprintf("0.0.0.0:%s", port),
			cert,
			key,
			nil,
		)
		if err != nil {
			log.Error(err, "error in http.ListenAndServeTLS")
		}
	}()

	return updates, nil
}

func getUpdatesChan(log *logr.Logger, tgAPI *tgbotapi.BotAPI) (updates tgbotapi.UpdatesChannel) {
	log.Info("starting tg_bot without webhook mode")

	updateCfg := tgbotapi.NewUpdate(0)
	updateCfg.Timeout = 60

	updates = tgAPI.GetUpdatesChan(updateCfg)
	return updates
}

func NewBot(
	ctx context.Context,
	debug, webhookMode bool,
	token, domain, port, cert, key string,
	svc pigeomail.Service,
	cons rabbitmq.Consumer,
) (bot *Bot, err error) {
	var log = logger.GetLogger()

	var tgAPI *tgbotapi.BotAPI
	if tgAPI, err = tgbotapi.NewBotAPI(token); err != nil {
		return nil, err
	}
	tgAPI.Debug = debug
	log.Info("authorized", "account", tgAPI.Self.UserName)

	log.Info("removing previous webhook")
	// delete created webhook cause
	// bot won't start in that mode if webhook was created before
	deleteWHCfg := tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: false,
	}

	if _, err = tgAPI.Request(deleteWHCfg); err != nil {
		log.Error(err, "fail remove webhook")
		return nil, err
	}

	var updates tgbotapi.UpdatesChannel
	switch webhookMode {
	case true:
		updates, err = getWebhookUpdatesChan(
			tgAPI,
			domain,
			port,
			cert,
			key,
		)
	case false:
		updates = getUpdatesChan(log, tgAPI)
	}
	// check get updates chan err
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:      tgAPI,
		updates:  updates,
		svc:      svc,
		consumer: cons,
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

	chatID, err := b.svc.GetChatIDByEmail(ctx, to.(string))
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
	err := b.consumer.Consume(rabbitmq.MessageReceivedQueueName, b.incomingEmailConsumer)
	if err != nil {
		b.logger.Error(err, "error runConsumer")
		os.Exit(1)
		return
	}
}

func (b *Bot) runBot() {
	for update := range b.updates {
		update := update
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(&update)
			continue
		}

		b.handleUserInput(&update)
	}
}

func (b *Bot) Run() {
	go b.runConsumer()
	b.runBot()
}

func (b *Bot) handleError(err error, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "")

	var customErr *customerrors.TelegramError
	if errors.As(err, &customErr) {
		msg.Text = customErr.Error()
		_, _ = b.api.Send(msg)
		return
	}

	uid := uuid.New().String()
	b.logger.Error(err, "unexpected error", "error_code", uid)
	msg.Text = "unexpected error, contact with support and send error code: " + uid
	_, _ = b.api.Send(msg)
}
