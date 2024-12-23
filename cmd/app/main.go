package main

import (
	"context"
	_ "go.uber.org/automaxprocs"

	"url-shortener/internal/app"
	"url-shortener/internal/config"
	"url-shortener/pkg/logger"
	"url-shortener/pkg/observability/otel"
	"url-shortener/pkg/observability/sentry"

	"github.com/rs/zerolog/log"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config.New")
	}

	logger.Init(c.Logger, c.App.Name, c.App.Version)
	log.Info().Msg("App starting...")
	defer log.Info().Msg("App stopped")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = sentry.Init(c.Sentry, c.App.Name, c.App.Version, c.App.Env)
	if err != nil {
		log.Error().Err(err).Msg("sentry.Init")
	}
	defer sentry.Close()

	err = otel.Init(ctx, c.Otel, c.App.Name, c.App.Version)
	if err != nil {
		log.Error().Err(err).Msg("otel.Init")
	}
	defer otel.Close()

	err = app.Run(ctx, c)
	if err != nil {
		log.Error().Err(err).Msg("app.Run")
	}
}
