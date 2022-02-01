package monitoring

import (
	"log"

	"github.com/getsentry/sentry-go"
)

func InitSentry(dsn string) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}
