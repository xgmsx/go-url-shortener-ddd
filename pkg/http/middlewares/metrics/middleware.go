package metrics

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/prometheus/client_golang/prometheus"
)

var HandlerDefault = NewHandler()

func NewHandler() fiber.Handler {
	gatherer := prometheus.DefaultGatherer
	if g, ok := registry.(prometheus.Gatherer); ok && g != gatherer {
		gatherer = g
	}

	return func(c *fiber.Ctx) error {
		return adaptor.HTTPHandler(promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}))(c)
	}
}

func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)
	if cfg.Registry == nil {
		cfg.Registry = prometheus.NewRegistry()
	}

	// Initializing metrics
	requestsTotal, requestsDuration, requestsInProgress := getMetrics(cfg.Registry, cfg.ServiceName, cfg.Labels)

	// Return middleware handler
	return func(c *fiber.Ctx) (err error) {
		// Пропустить, если функция Next() вернет true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Пропустить, если путь указан в списке исключений
		path := string(c.Request().RequestURI())
		if _, exists := cfg.skipPaths[path]; exists {
			return c.Next()
		}

		start := time.Now()
		method := utils.CopyString(c.Route().Method)

		// Изменение метрики http_requests_in_progress_total
		requestsInProgress.WithLabelValues(method, path).Inc()
		defer func() { requestsInProgress.WithLabelValues(method, path).Dec() }()

		// Выполнение запроса
		status := fiber.StatusInternalServerError
		if err = c.Next(); err == nil {
			status = c.Response().StatusCode()

			// Изменение метрики http_requests_duration если запрос успешный
			duration := float64(time.Since(start).Microseconds()) / 1e6
			requestsDuration.WithLabelValues(strconv.Itoa(status), method, path).Observe(duration)
		} else {
			if e, ok := err.(*fiber.Error); ok {
				status = e.Code
			}
		}

		// Изменение метрики http_requests_total
		requestsTotal.WithLabelValues(strconv.Itoa(status), method, path).Inc()

		return err
	}
}
