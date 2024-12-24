package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ErrHTTP struct {
	Error string `json:"error"`
}

func DefaultErrorHandler(c *fiber.Ctx, err error) error {
	if e, ok := err.(*fiber.Error); ok {
		return c.Status(e.Code).JSON(ErrHTTP{Error: e.Message})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(
		ErrHTTP{Error: http.StatusText(fiber.StatusInternalServerError)})
}
