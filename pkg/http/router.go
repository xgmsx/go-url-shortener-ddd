package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type DefaultRouter struct{}

func (r DefaultRouter) Register(prefix string, app *fiber.App) {
	router := app.Group(prefix)
	router.Get("/live", r.probe)
	router.Get("/ready", r.probe)
	router.Get("/swagger/*", swagger.HandlerDefault)
}

func (r DefaultRouter) probe(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
