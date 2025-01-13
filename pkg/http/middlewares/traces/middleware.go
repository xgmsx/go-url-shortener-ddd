package traces

import (
	"net/http"

	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

const (
	requestHeadersLen     = 16
	requestAttributesLen  = 16
	responseAttributesLen = 4
)

func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}

	// Initializing
	tracer := cfg.TracerProvider.Tracer("fiber",
		oteltrace.WithInstrumentationVersion(otelcontrib.Version()),
	)

	if cfg.Propagators == nil {
		cfg.Propagators = otel.GetTextMapPropagator()
	}
	if cfg.SpanNameFormatter == nil {
		cfg.SpanNameFormatter = defaultSpanNameFormatter
	}

	// Return new handler
	return func(c *fiber.Ctx) (err error) {
		// Пропустить, если функция Next() вернет true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Пропустить, если путь указан в списке исключений
		if _, exists := cfg.skipPaths[string(c.Request().RequestURI())]; exists {
			return c.Next()
		}

		// before request
		spanName := cfg.SpanNameFormatter(c)
		spanOpts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(getTraceAttributesFromRequest(c, &cfg)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}

		reqCtx := c.UserContext()
		defer c.SetUserContext(reqCtx)

		// extract TraceID from headers
		reqHeaders := make(http.Header, requestHeadersLen)
		c.Request().Header.VisitAll(func(k, v []byte) { reqHeaders.Add(string(k), string(v)) })
		ctx := cfg.Propagators.Extract(reqCtx, propagation.HeaderCarrier(reqHeaders))

		// create span
		spanCtx, span := tracer.Start(ctx, spanName, spanOpts...)
		defer span.End()
		c.SetUserContext(spanCtx)

		// request
		if err = c.Next(); err != nil {
			span.RecordError(err)
		}

		// inject TraceID into headers
		respHeaders := make(propagation.HeaderCarrier, 1)
		cfg.Propagators.Inject(spanCtx, respHeaders)
		for _, headerKey := range respHeaders.Keys() {
			c.Set(headerKey, respHeaders.Get(headerKey))
		}

		// after request
		span.SetAttributes(getTraceAttributesFromResponse(c)...)
		span.SetStatus(semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(
			c.Response().StatusCode(), oteltrace.SpanKindServer))

		return err
	}
}

func defaultSpanNameFormatter(c *fiber.Ctx) string {
	return string(utils.CopyBytes(c.Request().RequestURI()))
}

func getTraceAttributesFromRequest(c *fiber.Ctx, cfg *Config) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, requestAttributesLen)
	attrs = append(attrs,
		semconv.HTTPRequestContentLengthKey.Int(c.Request().Header.ContentLength()),
		semconv.HTTPMethodKey.String(utils.CopyString(c.Method())),
		semconv.HTTPSchemeKey.String(utils.CopyString(c.Protocol())),
		semconv.HTTPTargetKey.String(string(utils.CopyBytes(c.Request().RequestURI()))),
		semconv.HTTPURLKey.String(utils.CopyString(c.OriginalURL())),
		semconv.HTTPUserAgentKey.String(string(utils.CopyBytes(c.Request().Header.UserAgent()))),
		semconv.NetHostNameKey.String(utils.CopyString(c.Hostname())),
		semconv.NetTransportTCP,
	)

	if cfg.ServerName != "" {
		attrs = append(attrs, semconv.HTTPServerNameKey.String(cfg.ServerName))
	}

	if cfg.ServerPort != 0 {
		attrs = append(attrs, semconv.NetHostPortKey.Int(cfg.ServerPort))
	}

	if cfg.CollectClientIP {
		if clientIP := c.IP(); clientIP != "" {
			attrs = append(attrs, semconv.HTTPClientIPKey.String(utils.CopyString(clientIP)))
		}
	}

	return attrs
}

func getTraceAttributesFromResponse(c *fiber.Ctx) []attribute.KeyValue {
	var responseSize int64
	if c.GetRespHeader("Content-Type") != "text/event-stream" {
		responseSize = int64(len(c.Response().Body()))
	}

	attrs := make([]attribute.KeyValue, 0, responseAttributesLen)
	attrs = append(attrs, semconv.HTTPAttributesFromHTTPStatusCode(c.Response().StatusCode())...)
	attrs = append(attrs,
		semconv.HTTPResponseContentLengthKey.Int64(responseSize),
		semconv.HTTPRouteKey.String(c.Route().Path),
	)

	return attrs
}
