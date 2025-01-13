package traces

import (
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/gofiber/fiber/v2"
)

var ConfigDefault = Config{
	Next:              nil,
	SpanNameFormatter: nil,
	TracerProvider:    nil,
	Propagators:       nil,
	CollectClientIP:   true,
}

type Config struct {
	skipPaths         map[string]struct{}
	Next              func(*fiber.Ctx) bool
	SpanNameFormatter func(*fiber.Ctx) string
	TracerProvider    trace.TracerProvider
	Propagators       propagation.TextMapPropagator
	ServerName        string
	ServerPort        int
	CollectClientIP   bool
}

func (c *Config) SetSkipPaths(paths ...string) *Config {
	if len(paths) > 0 {
		c.skipPaths = make(map[string]struct{}, len(paths))
		for _, path := range paths {
			c.skipPaths[path] = struct{}{}
		}
	}
	return c
}

func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.skipPaths == nil {
		cfg.skipPaths = ConfigDefault.skipPaths
	}

	return cfg
}
