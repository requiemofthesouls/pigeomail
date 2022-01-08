package smtp_server

import (
	"log"
	"time"

	"github.com/emersion/go-smtp"
	"pigeomail/internal/repository"
	"pigeomail/rabbitmq"
)

type Receiver struct {
	server *smtp.Server
}

func NewSMTPReceiver(rmqCfg *rabbitmq.Config, cfg *Config, repo repository.IEmailRepository) (r *Receiver, err error) {
	var publisher rabbitmq.IRMQEmailPublisher
	if publisher, err = rabbitmq.NewRMQEmailPublisher(rmqCfg); err != nil {
		return nil, err
	}

	var b smtp.Backend
	if b, err = NewBackend(publisher, repo); err != nil {
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

	return &Receiver{server: server}, nil
}

func (r *Receiver) Run() (err error) {
	log.Println("Starting Receiver at", r.server.Addr)
	if err = r.server.ListenAndServe(); err != nil {
		return err
	}
	return err
}
