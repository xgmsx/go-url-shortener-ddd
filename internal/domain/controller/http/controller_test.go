package http

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/create"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/fetch"
)

func TestController(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error { return c.SendString("Test response") })

	ctrl := New("/test/api", create.Usecase{}, fetch.Usecase{})
	ctrl.Register(app)

	resp, body := sendHTTPRequest(t, app, http.MethodGet, "/test", "")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Test response", body)

	resp, body = sendHTTPRequest(t, app, http.MethodGet, "/test/api/link", "")
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	assert.Equal(t, "Method Not Allowed", body)

	resp, body = sendHTTPRequest(t, app, http.MethodGet, "/unknown", "")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "Cannot GET /unknown", body)
}
