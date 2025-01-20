package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/usecase/create"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/usecase/fetch"
)

type Router struct {
	prefix   string
	ucCreate create.Usecase
	ucFetch  fetch.Usecase
}

func New(prefix string, ucCreate create.Usecase, ucFetch fetch.Usecase) *Router {
	return &Router{
		prefix:   prefix,
		ucCreate: ucCreate,
		ucFetch:  ucFetch,
	}
}

func (r *Router) Register(app *fiber.App) {
	router := app.Group(r.prefix)

	v1 := router.Group("/v1")
	v1.Post("/link", r.createLink)
	v1.Get("/link/:alias", r.fetchLink)
	v1.Get("/link/:alias/redirect", r.redirect)
}
