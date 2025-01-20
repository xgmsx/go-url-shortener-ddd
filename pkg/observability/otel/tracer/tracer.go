package tracer

import (
	"context"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func init() { // nolint // default tracer
	tracer = otel.Tracer("")
}

func Init(t trace.Tracer) {
	tracer = t
}

func Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName, opts...)
}

func SetStatus(span trace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
}
