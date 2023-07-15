package telegram

import (
	"context"
	"fmt"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
	"go.uber.org/zap"
)

type (
	UserInputHandler   func(update *tgbotapi.Update)
	UserCommandHandler func(update *tgbotapi.Update)

	Wrapper interface {
		Start(
			ctx context.Context,
			userInputHandler UserInputHandler,
			userCommandHandler UserCommandHandler,
		)

		Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	}

	wrapper struct {
		api     *tgbotapi.BotAPI
		updates tgbotapi.UpdatesChannel
	}
)

func (w *wrapper) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return w.api.Send(c)
}

func New(
	cfg *Config,
	l logger.Wrapper,
) (Wrapper, error) {
	var (
		tgAPI *tgbotapi.BotAPI
		err   error
	)
	if tgAPI, err = tgbotapi.NewBotAPI(cfg.Token); err != nil {
		return nil, fmt.Errorf("tgbotapi.NewBotAPI error: %w", err)
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
		return nil, fmt.Errorf("tgAPI.Request error: %w", err)
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
		updates = getUpdatesChan(
			l,
			tgAPI,
		)
	}
	// check get updates chan err
	if err != nil {
		return nil, err
	}

	return &wrapper{
		api:     tgAPI,
		updates: updates,
	}, nil
}

func getUpdatesChan(
	l logger.Wrapper,
	tgAPI *tgbotapi.BotAPI,
) (updates tgbotapi.UpdatesChannel) {
	l.Info("starting tg_bot without webhook mode")

	updateCfg := tgbotapi.NewUpdate(0)
	updateCfg.Timeout = 60

	updates = tgAPI.GetUpdatesChan(updateCfg)
	return updates
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
		server := &http.Server{
			Addr:              fmt.Sprintf("0.0.0.0:%d", cfg.Port),
			ReadHeaderTimeout: 3 * time.Second,
		}

		if err = server.ListenAndServeTLS(cfg.Cert, cfg.Key); err != nil {
			l.Error("error in http.ListenAndServeTLS", zap.Error(err))
		}
	}()

	return updates, nil
}

func (w *wrapper) runBot(
	userInputHandler UserInputHandler,
	userCommandHandler UserCommandHandler,
) {
	for update := range w.updates {
		update := update
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			userCommandHandler(&update)
			continue
		}

		userInputHandler(&update)
	}
}

func (w *wrapper) Start(
	ctx context.Context,
	userInputHandler UserInputHandler,
	userCommandHandler UserCommandHandler,
) {
	l := getLogger(ctx)

	l.Info("starting tg_bot")
	go w.runBot(userInputHandler, userCommandHandler)
	<-ctx.Done()
	l.Info("stopping tg_bot")
}

func getLogger(ctx context.Context) logger.Wrapper {
	return logger.NewFromZap(ctxzap.Extract(ctx))
}
