package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/create"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/fetch"
)

type Controller struct {
	prefix   string
	ucCreate create.Usecase
	ucFetch  fetch.Usecase
}

func New(prefix string, ucCreate create.Usecase, ucFetch fetch.Usecase) *Controller {
	return &Controller{prefix, ucCreate, ucFetch}
}

func (c *Controller) Register(app *fiber.App) {
	r := app.Group(c.prefix)
	r.Post("/link", NewHandlerCreateLink(c.ucCreate).Handler)
	r.Get("/link/:alias", NewHandlerFetchLink(c.ucFetch).Handler)
	r.Get("/link/:alias/redirect", NewHandlerRedirect(c.ucFetch).Handler)
}
