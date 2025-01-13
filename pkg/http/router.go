package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"github.com/xgmsx/go-url-shortener-ddd/pkg/http/middlewares/metrics"
)

type DefaultRouter struct{}

func (r DefaultRouter) Register(prefix string, app *fiber.App) {
	router := app.Group(prefix)
	router.Get("/live", r.probe)
	router.Get("/ready", r.probe)
	router.Get("/swagger/*", swagger.HandlerDefault)

	metrics.RegisterAt(app, "/metrics")
}

func (DefaultRouter) probe(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
