package receiver

import (
	"time"

	"github.com/emersion/go-smtp"
	"github.com/go-logr/logr"

	"pigeomail/pkg/logger"
)

type Receiver struct {
	server *smtp.Server
	logger *logr.Logger
}

func NewSMTPReceiver(
	backend smtp.Backend,
	addr, domain string,
	readTimeout, writeTimeout, maxMessageBytes, maxRecipients int,
	allowInsecureAuth bool,
) (r *Receiver, err error) {
	var log = logger.GetLogger()

	server := smtp.NewServer(backend)
	server.Addr = addr
	server.Domain = domain
	server.ReadTimeout = time.Duration(readTimeout) * time.Second
	server.WriteTimeout = time.Duration(writeTimeout) * time.Second
	server.MaxMessageBytes = maxMessageBytes * 1024
	server.MaxRecipients = maxRecipients
	server.AllowInsecureAuth = allowInsecureAuth

	return &Receiver{server: server, logger: log}, nil
}

func (r *Receiver) Run() (err error) {
	r.logger.Info("starting receiver", "addr", r.server.Addr)
	if err = r.server.ListenAndServe(); err != nil {
		return err
	}
	return err
}
