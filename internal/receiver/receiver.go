package receiver

import (
	"context"
	"time"

	"github.com/emersion/go-smtp"
	"go.uber.org/zap"

	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
)

type Receiver struct {
	server *smtp.Server
	logger logger.Wrapper
}

func NewSMTPReceiver(
	backend smtp.Backend,
	logger logger.Wrapper,
	cfg ServerConfig,
) (r *Receiver, err error) {
	server := smtp.NewServer(backend)
	server.Addr = cfg.getAddr()
	server.Domain = cfg.Domain
	server.ReadTimeout = time.Duration(cfg.ReadTimeoutSeconds) * time.Second
	server.WriteTimeout = time.Duration(cfg.WriteTimeoutSeconds) * time.Second
	server.MaxMessageBytes = cfg.MaxMessageBytes * 1024
	server.MaxRecipients = cfg.MaxRecipients
	server.AllowInsecureAuth = cfg.AllowInsecureAuth

	return &Receiver{server: server, logger: logger}, nil
}

func (r *Receiver) Run(ctx context.Context) {
	r.logger.Info(
		"starting receiver",
		zap.String("addr", r.server.Addr),
		zap.String("domain", r.server.Domain),
	)

	go func() {
		if err := r.server.ListenAndServe(); err != nil {
			r.logger.Error("ListenAndServe error", zap.Error(err))
			panic(err)
		}
	}()

	<-ctx.Done()

	r.logger.Info("stopping receiver")
}
