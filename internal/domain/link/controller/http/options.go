package http

import (
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/xgmsx/go-url-shortener-ddd/internal/config"
	hs "github.com/xgmsx/go-url-shortener-ddd/pkg/http"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/http/middlewares/metrics"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/http/middlewares/traces"
)

func DefaultOptions(c *config.Config) *hs.Options {
	mc := metrics.Config{ServiceName: c.App.Name, Registry: prometheus.DefaultRegisterer}
	tc := traces.Config{ServerName: c.App.Name, CollectClientIP: true}
	tc.SetSkipPaths("/metrics", "/live", "/ready")

	return hs.BuildOptions(c.HTTP,
		hs.WithMiddleware(requestid.New(requestid.Config{}), c.HTTP.UseRequestID),
		hs.WithMiddleware(pprof.New(pprof.Config{}), c.HTTP.UsePprof),
		hs.WithMiddleware(traces.New(tc), c.Otel.Endpoint != ""),
		hs.WithMiddleware(metrics.New(mc)),
		hs.WithMiddleware(recover.New(recover.Config{}), c.HTTP.UseRecover),
	)
}
