package http_router

import (
	"url-shortener/internal/shortener/usecase"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	uc *usecase.UseCase
}

func New(uc *usecase.UseCase) Router {
	return Router{uc: uc}
}

func (r Router) Register(prefix string, app *fiber.App) {
	router := app.Group(prefix)

	v1 := router.Group("/v1")
	v1.Post("/link", r.createLink)
	v1.Get("/link/:alias", r.getLink)
	v1.Get("/link/:alias/redirect", r.redirect)
}
