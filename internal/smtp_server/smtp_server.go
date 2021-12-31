package smtp_server

import (
	"log"

	"github.com/emersion/go-smtp"
)

func Run(s *smtp.Server) (err error) {
	log.Println("Starting server at", s.Addr)
	if err = s.ListenAndServe(); err != nil {
		return err
	}
	return err
}
