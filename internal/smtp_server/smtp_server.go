package smtp_server

import (
	"time"

	"github.com/emersion/go-smtp"
	"github.com/go-logr/logr"
	"pigeomail/internal/repository"
	"pigeomail/rabbitmq"
)

type Receiver struct {
	server *smtp.Server
	logger logr.Logger
}

func NewSMTPReceiver(rmqCfg *rabbitmq.Config, cfg *Config, repo repository.IEmailRepository, log logr.Logger) (r *Receiver, err error) {
	var publisher rabbitmq.IRMQEmailPublisher
	if publisher, err = rabbitmq.NewRMQEmailPublisher(rmqCfg); err != nil {
		return nil, err
	}

	var b smtp.Backend
	if b, err = NewBackend(publisher, repo, log); err != nil {
		return nil, err
	}

	server := smtp.NewServer(b)

	server.Addr = cfg.Addr
	server.Domain = cfg.Domain
	server.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Second
	server.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Second
	server.MaxMessageBytes = cfg.MaxMessageBytes * 1024
	server.MaxRecipients = cfg.MaxRecipients
	server.AllowInsecureAuth = cfg.AllowInsecureAuth

	return &Receiver{server: server, logger: log}, nil
}

func (r *Receiver) Run() (err error) {
	r.logger.Info("starting receiver", "addr", r.server.Addr)
	if err = r.server.ListenAndServe(); err != nil {
		return err
	}
	return err
}
