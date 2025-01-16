package main

import (
	"context"

	"github.com/rs/zerolog/log"
	_ "go.uber.org/automaxprocs"

	"github.com/xgmsx/go-url-shortener-ddd/internal/app"
	"github.com/xgmsx/go-url-shortener-ddd/internal/config"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/logger"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/sentry"
)

func run(run func(context.Context, *config.Config) error) {
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

	err = run(ctx, c)
	if err != nil {
		log.Error().Err(err).Msg("app.Run")
	}
}

func main() {
	run(app.Run)
}
