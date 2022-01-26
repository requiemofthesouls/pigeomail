package monitoring

import (
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
)

func InitSentry() {

	// TODO: Need help in unmarshalling this config properly
	// now it is crashing

	var cfg *Config
	if err := viper.UnmarshalKey("monitoring", &cfg); err != nil {
		log.Println("Can't unmarshal config of sentry")
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn: cfg.SentryDsn,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}
