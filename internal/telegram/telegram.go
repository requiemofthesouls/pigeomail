package telegram

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"github.com/requiemofthesouls/pigeomail/internal/customerrors"
	"github.com/requiemofthesouls/pigeomail/internal/rabbitmq"
	"github.com/requiemofthesouls/pigeomail/internal/repository"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	updates  tgbotapi.UpdatesChannel
	repo     repository.EmailState
	consumer rabbitmq.Consumer
	domain   string
	logger   logger.Wrapper
}

func getWebhookUpdatesChan(
	tgAPI *tgbotapi.BotAPI,
	l logger.Wrapper,
	cfg WebhookConfig,
) (updates tgbotapi.UpdatesChannel, err error) {
	l.Info("starting tg_bot in webhook mode", zap.Uint32("port", cfg.Port))

	var whCfg tgbotapi.WebhookConfig
	whCfg, err = tgbotapi.NewWebhookWithCert(
		fmt.Sprintf("https://%s:%d/%s", cfg.Domain, cfg.Port, tgAPI.Token),
		tgbotapi.FilePath(cfg.Cert),
	)
	if err != nil {
		l.Error("fail to initialize tgbotapi.NewWebhookWithCert", zap.Error(err))
		return nil, err
	}

	if _, err = tgAPI.Request(whCfg); err != nil {
		l.Error("fail set webhook", zap.Error(err))
		return nil, err
	}

	var info tgbotapi.WebhookInfo
	if info, err = tgAPI.GetWebhookInfo(); err != nil {
		l.Error("fail GetWebhookInfo", zap.Error(err))
		return nil, err
	}

	if info.LastErrorDate != 0 {
		l.Info("GetWebhookInfo", zap.String("last_error", info.LastErrorMessage))
	}

	updates = tgAPI.ListenForWebhook("/" + tgAPI.Token)

	go func() {
		err = http.ListenAndServeTLS(
			fmt.Sprintf("0.0.0.0:%d", cfg.Port),
			cfg.Cert,
			cfg.Key,
			nil,
		)
		if err != nil {
			l.Error("error in http.ListenAndServeTLS", zap.Error(err))
		}
	}()

	return updates, nil
}

func getUpdatesChan(l logger.Wrapper, tgAPI *tgbotapi.BotAPI) (updates tgbotapi.UpdatesChannel) {
	l.Info("starting tg_bot without webhook mode")

	updateCfg := tgbotapi.NewUpdate(0)
	updateCfg.Timeout = 60

	updates = tgAPI.GetUpdatesChan(updateCfg)
	return updates
}

func NewBot(
	cfg Config,
	l logger.Wrapper,
	repo repository.EmailState,
	cons rabbitmq.Consumer,
) (bot *Bot, err error) {
	var tgAPI *tgbotapi.BotAPI
	if tgAPI, err = tgbotapi.NewBotAPI(cfg.Token); err != nil {
		return nil, err
	}
	tgAPI.Debug = cfg.Debug
	l.Info("authorized", zap.String("account", tgAPI.Self.UserName))

	l.Info("removing previous webhook")
	// delete created webhook cause
	// bot won't start in that mode if webhook was created before
	deleteWHCfg := tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: false,
	}

	if _, err = tgAPI.Request(deleteWHCfg); err != nil {
		l.Error("fail remove webhook", zap.Error(err))
		return nil, err
	}

	var updates tgbotapi.UpdatesChannel
	switch cfg.Webhook.Enabled {
	case true:
		updates, err = getWebhookUpdatesChan(
			tgAPI,
			l,
			cfg.Webhook,
		)
	case false:
		updates = getUpdatesChan(l, tgAPI)
	}
	// check get updates chan err
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:      tgAPI,
		updates:  updates,
		repo:     repo,
		consumer: cons,
		domain:   cfg.Webhook.Domain,
		logger:   l,
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
			b.logger.Error("error send message", zap.Error(err))
		}
	}
}

func (b *Bot) incomingEmailConsumer(msg *amqp.Delivery) {
	from, ok := msg.Headers["from"]
	if !ok {
		b.logger.Warn("fail to extract 'from' header from message")
		_ = msg.Reject(false)
	}

	to, ok := msg.Headers["to"]
	if !ok {
		b.logger.Warn("fail to extract 'to' header from message")
		_ = msg.Reject(false)
	}

	subject, ok := msg.Headers["subject"]
	if !ok {
		b.logger.Warn("fail to extract 'subject' header from message")
		_ = msg.Reject(false)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chatID, err := b.repo.GetChatIDByEmail(ctx, to.(string))
	if err != nil {
		b.logger.Error("chatID not found", zap.String("email", to.(string)), zap.Error(err))
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
			b.logger.Error("error send message", zap.Error(err))
		}

		for i := 3000; i < len(msg.Body); i += 4096 {
			y := i + 4096
			if y > len(msg.Body) {
				y = len(msg.Body)
			}

			tgMsg = tgbotapi.NewMessage(chatID, html.EscapeString(string(msg.Body[i:y])))
			tgMsg.ParseMode = tgbotapi.ModeHTML

			if _, err = b.api.Send(tgMsg); err != nil {
				b.logger.Error("error send message", zap.Error(err))
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
		b.logger.Error("error send message", zap.Error(err))
	}

	_ = msg.Ack(false)
}

func (b *Bot) runConsumer() {
	err := b.consumer.Consume(rabbitmq.MessageReceivedQueueName, b.incomingEmailConsumer)
	if err != nil {
		b.logger.Error("error runConsumer", zap.Error(err))
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
	b.logger.Error("unexpected error", zap.String("error_code", uid), zap.Error(err))
	msg.Text = "unexpected error, contact with support and send error code: " + uid
	_, _ = b.api.Send(msg)
}
