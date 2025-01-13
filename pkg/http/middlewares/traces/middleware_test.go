package traces

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func getTestOtelSpanRecord() *tracetest.SpanRecorder {
	spanRecorder := tracetest.NewSpanRecorder()
	traceProvider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return spanRecorder
}

func getBenchRequest(path string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	req := &fasthttp.Request{}
	req.Header.SetMethod(fiber.MethodOptions)
	req.SetRequestURI(path)
	ctx.Init(req, nil, nil)

	return ctx
}

func TestTracesMiddleware(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		path    string
		wantLen int
	}{
		{
			name:    "With default config",
			config:  ConfigDefault,
			path:    "/health",
			wantLen: 1,
		},
		{
			name:    "With custom config",
			config:  Config{ServerName: "Test", ServerPort: 8080, CollectClientIP: true},
			path:    "/health",
			wantLen: 1,
		},
		{
			name:    "With Next function",
			config:  Config{Next: func(c *fiber.Ctx) bool { return c.Path() == "/health" }},
			path:    "/health",
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			sr := getTestOtelSpanRecord()

			app := fiber.New()
			app.Use(New(tt.config))
			app.Get("/health", func(ctx *fiber.Ctx) error { return ctx.SendString("ok") })

			// Act
			_, _ = app.Test(httptest.NewRequest("GET", tt.path, http.NoBody))

			// Assert
			assert.Len(t, sr.Ended(), tt.wantLen)
		})
	}
}

func TestTracesMiddlewareWithError(t *testing.T) {
	// Arrange
	sr := getTestOtelSpanRecord()

	app := fiber.New()
	app.Use(New())
	app.Get("/health", func(ctx *fiber.Ctx) error { return ctx.SendString("ok") })

	// Act
	_, _ = app.Test(httptest.NewRequest("GET", "/unknown", http.NoBody))

	// Assert
	assert.Len(t, sr.Ended(), 1)
}

func TestTracesMiddlewarePropagationInject(t *testing.T) {
	// Arrange
	sr := getTestOtelSpanRecord()

	app := fiber.New()
	app.Use(New())
	app.Get("/health", func(ctx *fiber.Ctx) error { return ctx.SendString("ok") })

	// Act
	resp, _ := app.Test(httptest.NewRequest("GET", "/health", http.NoBody))
	require.Len(t, sr.Ended(), 1)

	// Assert
	traceparentHeader := resp.Header.Get("traceparent")
	assert.NotEmpty(t, traceparentHeader)

	spanTraceID := sr.Ended()[0].SpanContext().TraceID().String()
	assert.Contains(t, traceparentHeader, spanTraceID)
}

func TestTracesMiddlewareMemory(t *testing.T) {
	// Arrange
	sr := getTestOtelSpanRecord()

	app := fiber.New()
	app.Use(New())
	app.Get("/health", func(ctx *fiber.Ctx) error { return ctx.SendStatus(201) })
	app.Get("/ping", func(ctx *fiber.Ctx) error { return ctx.SendString("pong") })

	// Act
	_, _ = app.Test(httptest.NewRequest("GET", "/ping", http.NoBody))
	_, _ = app.Test(httptest.NewRequest("GET", "/health", http.NoBody))

	// Assert
	spans := sr.Ended()
	assert.Len(t, spans, 2)
	assert.Equal(t, "/ping", spans[0].Name())
	assert.Equal(t, "/health", spans[1].Name())
}

func TestTracesMiddlewarePropagationExtract(t *testing.T) {
	const traceparent = "00-c1c00fe240f3daa80eb223f96cc16f9f-884207f0a4f22360-01"

	// Arrange
	sr := getTestOtelSpanRecord()

	app := fiber.New()
	app.Use(New())
	app.Get("/health", func(ctx *fiber.Ctx) error { return ctx.SendString("ok") })

	// Act
	req := httptest.NewRequest("GET", "/health", http.NoBody)
	req.Header.Set("traceparent", traceparent)

	resp, _ := app.Test(req)
	require.Len(t, sr.Ended(), 1)

	// Assert
	parentTraceID := sr.Ended()[0].SpanContext().TraceID().String()
	assert.Contains(t, traceparent, parentTraceID)

	traceparentHeader := resp.Header.Get("Traceparent")
	assert.NotEmpty(t, traceparentHeader)

	childSpanID := sr.Ended()[0].SpanContext().SpanID().String()
	assert.Contains(t, traceparentHeader, parentTraceID)
	assert.Contains(t, traceparentHeader, childSpanID)
}

func TestTracesMiddlewareWithSkipPath(t *testing.T) {
	// Arrange
	sr := getTestOtelSpanRecord()

	cfg := Config{}
	cfg.SetSkipPaths("/test2")

	app := fiber.New()
	app.Use(New(cfg))
	app.Get("/test1", func(c *fiber.Ctx) error { return c.SendString("Test 1") })
	app.Get("/test2", func(c *fiber.Ctx) error { return c.SendString("Test 2") })

	// Act
	_, _ = app.Test(httptest.NewRequest("GET", "/test1", http.NoBody), -1)
	_, _ = app.Test(httptest.NewRequest("GET", "/test2", http.NoBody), -1)

	// Assert
	spans := sr.Ended()
	assert.Len(t, spans, 1)
}

func BenchmarkTracesMiddleware(b *testing.B) {
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
