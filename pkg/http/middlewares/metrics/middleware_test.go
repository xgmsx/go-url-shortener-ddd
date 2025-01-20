package metrics

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func getMetricsStr(app *fiber.App) string {
	resp, _ := app.Test(httptest.NewRequest("GET", "/metrics", http.NoBody), -1)
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func getBenchRequest(path string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	req := &fasthttp.Request{}
	req.Header.SetMethod(fiber.MethodOptions)
	req.SetRequestURI(path)
	ctx.Init(req, nil, nil)

	return ctx
}

func TestMetricsMiddlewareWithLabels(t *testing.T) {
	// Arrange
	app := fiber.New()
	app.Use(New(Config{ServiceName: "test service", Labels: map[string]string{"env": "preprod", "version": "1.0.0"}}))
	app.Get("/metrics", NewHandler())
	app.Get("/test", func(c *fiber.Ctx) error { return c.SendString("Test metrics") })

	// Act
	resp, _ := app.Test(httptest.NewRequest("GET", "/test", http.NoBody), -1)
	require.Equal(t, 200, resp.StatusCode)

	// Assert
	got := getMetricsStr(app)
	want := `http_requests_total{env="preprod",method="GET",path="/test",service="test service",status="200",version="1.0.0"} 1`
	assert.Contains(t, got, want)
}

func TestMetricsMiddlewareWithNext(t *testing.T) {
	// Arrange
	app := fiber.New()
	app.Use(New(Config{Next: func(c *fiber.Ctx) bool { return c.Path() == "/test2" }}))
	app.Get("/metrics", NewHandler())
	app.Get("/test1", func(c *fiber.Ctx) error { return c.SendString("Test 1") })
	app.Get("/test2", func(c *fiber.Ctx) error { return c.SendString("Test 2") })

	// Act
	resp, _ := app.Test(httptest.NewRequest("GET", "/test1", http.NoBody), -1)
	require.Equal(t, 200, resp.StatusCode)

	resp, _ = app.Test(httptest.NewRequest("GET", "/test2", http.NoBody), -1)
	require.Equal(t, 200, resp.StatusCode)

	// Assert
	got := getMetricsStr(app)
	assert.Contains(t, got, "/test1")
	assert.NotContains(t, got, "/test2")
}

func TestMetricsMiddlewareWithSkipPath(t *testing.T) {
	// Arrange
	cfg := Config{}
	cfg.SetSkipPaths("/test2", "/metrics")

	app := fiber.New()
	app.Use(New(cfg))
	app.Get("/metrics", NewHandler())
	app.Get("/test1", func(c *fiber.Ctx) error { return c.SendString("Test 1") })
	app.Get("/test2", func(c *fiber.Ctx) error { return c.SendString("Test 2") })

	// Act
	resp, _ := app.Test(httptest.NewRequest("GET", "/test1", http.NoBody), -1)
	require.Equal(t, 200, resp.StatusCode)

	resp, _ = app.Test(httptest.NewRequest("GET", "/test2", http.NoBody), -1)
	require.Equal(t, 200, resp.StatusCode)

	// Assert
	got := getMetricsStr(app)
	assert.Contains(t, got, "/test1")
	assert.NotContains(t, got, "/test2")
}

func TestMetricsMiddlewareWithBasicAuth(t *testing.T) {
	// Arrange
	app := fiber.New()
	app.Use(New())
	app.Get("/test", func(c *fiber.Ctx) error { return c.SendString("Hello World") })

	// Act
	app.Get("/metrics", basicauth.New(basicauth.Config{Users: map[string]string{"login": "pass"}}), NewHandler())

	// Assert
	resp, _ := app.Test(httptest.NewRequest("GET", "/test", http.NoBody), -1)
	require.Equal(t, 200, resp.StatusCode)

	req := httptest.NewRequest("GET", "/metrics", http.NoBody)
	resp, _ = app.Test(req, -1)
	require.Equal(t, 401, resp.StatusCode)

	req.SetBasicAuth("wrong_login", "wrong_pass")
	resp, _ = app.Test(req, -1)
	assert.Equal(t, 401, resp.StatusCode)

	req.SetBasicAuth("login", "pass")
	resp, _ = app.Test(req, -1)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestMetricsMiddlewareWithError(t *testing.T) {
	// Arrange
	app := fiber.New()
	app.Use(New())
	app.Get("/metrics", NewHandler())

	// Act
	resp, _ := app.Test(httptest.NewRequest("POST", "/unknown", http.NoBody), -1)
	require.Equal(t, 404, resp.StatusCode)

	// Assert
	got := getMetricsStr(app)
	want := `http_requests_total{method="POST",path="/unknown",status="404"} 1`
	assert.Contains(t, got, want)
}

func BenchmarkMetricsMiddleware(b *testing.B) {
	app := fiber.New()
	app.Use(New())
	app.Get("/bench", func(c *fiber.Ctx) error { return c.SendStatus(200) })

	handler := app.Handler()
	req := getBenchRequest("/bench")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handler(req)
	}
}
