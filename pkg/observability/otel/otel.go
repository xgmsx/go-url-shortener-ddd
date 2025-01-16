package otel

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	tracer_noop "go.opentelemetry.io/otel/trace/noop"

	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
)

const closeTimeout = 5 * time.Second

type Config struct {
	Endpoint     string  `env:"OTEL_ENDPOINT"`
	EndpointHTTP string  `env:"OTEL_ENDPOINT_HTTP"`
	Namespace    string  `env:"OTEL_NAMESPACE"`
	InstanceID   string  `env:"OTEL_INSTANCE_ID"`
	Ratio        float64 `env:"OTEL_RATIO, default=1.0"`
}

var shutdownTracing func(ctx context.Context) error

func Init(ctx context.Context, c Config, name, version string) error {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

	var exporter *otlptrace.Exporter
	var err error

	switch {
	case c.Endpoint != "":
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(c.Endpoint),
			otlptracegrpc.WithInsecure())
		if err != nil {
			return fmt.Errorf("failed to create OTLP trace exporter (GRPC): %w", err)
		}
	case c.EndpointHTTP != "":
		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(c.EndpointHTTP),
			otlptracehttp.WithInsecure())
		if err != nil {
			return fmt.Errorf("failed to create OTLP trace exporter (HTTP): %w", err)
		}
	default:
		otel.SetTracerProvider(tracer_noop.NewTracerProvider())
		tracer.Init(otel.Tracer(""))
		log.Info().Msg("Tracer is disabled")
		return nil
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
			semconv.ServiceVersionKey.String(version),
			semconv.ServiceNamespaceKey.String(c.Namespace),
			semconv.ServiceInstanceIDKey.String(c.InstanceID),
		)),
	)

	shutdownTracing = traceProvider.Shutdown
	otel.SetTracerProvider(traceProvider)
	tracer.Init(otel.Tracer("app"))

	log.Info().Msg("Tracer initialized")
	return nil
}

func Close() {
	if shutdownTracing == nil {
		return
	}

	ctx, stop := context.WithTimeout(context.Background(), closeTimeout)
	defer stop()

	err := shutdownTracing(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to shutdown tracing")
	}

	log.Info().Msg("Tracer closed")
}
