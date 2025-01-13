package metrics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var ConfigDefault = Config{
	skipPaths:       defaultSkipPaths,
	Labels:          nil,
	Next:            nil,
	Registry:        nil,
	DurationBuckets: defaultBuckets,
	ServiceName:     "",
}

type Config struct {
	skipPaths       map[string]struct{}
	Labels          map[string]string
	Next            func(c *fiber.Ctx) bool
	Registry        prometheus.Registerer
	DurationBuckets []float64
	ServiceName     string
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
	if cfg.DurationBuckets == nil {
		cfg.DurationBuckets = ConfigDefault.DurationBuckets
	}
	if cfg.Registry == nil {
		cfg.Registry = ConfigDefault.Registry
	}
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}

	return cfg
}
