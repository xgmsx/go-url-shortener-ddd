package config

import (
	"context"

	"url-shortener/pkg/grpc"
	"url-shortener/pkg/http"
	"url-shortener/pkg/kafka/reader"
	"url-shortener/pkg/kafka/writer"
	"url-shortener/pkg/logger"
	"url-shortener/pkg/observability/otel"
	"url-shortener/pkg/observability/sentry"
	"url-shortener/pkg/postgres"
	"url-shortener/pkg/redis"

	"github.com/sethvargo/go-envconfig"
)

type App struct {
	Name    string `env:"APP_NAME, default=url-shortener"`
	Version string `env:"APP_VERSION, default=0.0.0"`
	Env     string `env:"APP_ENV, default=DEV"`
}

type Config struct {
	App    App
	Logger logger.Config
	// Observability
	Sentry sentry.Config
	Otel   otel.Config
	// Dependencies
	Postgres    postgres.Config
	Redis       redis.Config
	KafkaWriter writer.Config
	KafkaReader reader.Config
	// Controllers
	HTTP http.Config
	GRPC grpc.Config
}

func New() (Config, error) {
	var cfg Config
	var err = envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
