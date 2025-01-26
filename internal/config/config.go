package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"

	"github.com/xgmsx/go-url-shortener-ddd/pkg/grpc"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/http"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/kafka/reader"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/kafka/writer"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/logger"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/sentry"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/postgres"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/redis"
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

func New() *Config {
	return &Config{}
}

func (c *Config) Load(ctx context.Context) (*Config, error) {
	err := envconfig.Process(ctx, c)
	return c, err
}
