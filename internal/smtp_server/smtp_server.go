package smtp_server

import (
	"log"
	"time"

	"github.com/emersion/go-smtp"
	"pigeomail/config"
)

// build Builds smtp server with options in config
func build() (s *smtp.Server, err error) {
	var c *config.Config
	if c, err = config.Get(); err != nil {
		return nil, err
	}

	var b smtp.Backend
	if b, err = NewBackend(); err != nil {
		return nil, err
	}

	s = smtp.NewServer(b)

	s.Addr = c.Addr
	s.Domain = c.Domain
	s.ReadTimeout = time.Duration(c.ReadTimeout) * time.Second
	s.WriteTimeout = time.Duration(c.WriteTimeout) * time.Second
	s.MaxMessageBytes = c.MaxMessageBytes * 1024
	s.MaxRecipients = c.MaxRecipients
	s.AllowInsecureAuth = c.AllowInsecureAuth

	return s, nil
}

func Run() (err error) {
	var s *smtp.Server
	if s, err = build(); err != nil {
		return err
	}

	log.Println("Starting server at", s.Addr)
	if err = s.ListenAndServe(); err != nil {
		return err
	}

	return err
}
