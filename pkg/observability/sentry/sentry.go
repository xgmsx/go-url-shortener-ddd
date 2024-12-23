package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
)

const closeTimeout = 5 * time.Second

type Config struct {
	DSN              string  `env:"SENTRY_DSN"`
	Rate             float64 `env:"SENTRY_RATE, default=1.0"`
	AttachStackTrace bool    `env:"SENTRY_STACK_TRACE, default=true"`
	Debug            bool    `env:"SENTRY_DEBUG, default=false"`
}

func Init(c Config, name, version, env string) error {
	if c.DSN == "" {
		log.Info().Msg("Sentry is disabled")
		return nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              c.DSN,
		SampleRate:       c.Rate,
		AttachStacktrace: c.AttachStackTrace,
		Debug:            c.Debug,
		ServerName:       name,
		Release:          version,
		Environment:      env,
	})
	if err != nil {
		return err
	}

	log.Info().Msg("Sentry initialized")
	return nil
}

func Close() {
	_ = sentry.Flush(closeTimeout)
	log.Info().Msg("Sentry closed")
}
