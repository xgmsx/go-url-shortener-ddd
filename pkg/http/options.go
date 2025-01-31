package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/xgmsx/go-url-shortener-ddd/pkg/http/middlewares/metrics"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/http/middlewares/traces"
)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 5 * time.Second
	closeTimeout = 5 * time.Second
)

type Option func(*Options)

type Options struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	CloseTimeout time.Duration
	Middlewares  []fiber.Handler
}

func defaultOptions(c Config, options ...*Options) *Options {
	if options != nil && options[0] != nil {
		return options[0]
	}

	mc := metrics.Config{ServiceName: c.AppName, Registry: prometheus.DefaultRegisterer}
	tc := traces.Config{ServerName: c.AppName, CollectClientIP: true}
	tc.SetSkipPaths("/metrics", "/live", "/ready")

	o := &Options{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		CloseTimeout: closeTimeout,
	}

	for _, opt := range []Option{
		WithMiddleware(requestid.New(requestid.Config{}), c.UseRequestID),
		WithMiddleware(pprof.New(pprof.Config{}), c.UsePprof),
		WithMiddleware(traces.New(tc)),
		WithMiddleware(metrics.New(mc)),
		WithMiddleware(recover.New(recover.Config{}), c.UseRecover),
	} {
		opt(o)
	}

	return o
}

func WithMiddleware(m fiber.Handler, conditions ...bool) Option {
	return func(options *Options) {
		for _, v := range conditions {
			if !v {
				return
			}
		}

		options.Middlewares = append(options.Middlewares, m)
	}
}
