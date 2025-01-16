package kafka_producer //nolint:stylecheck

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/entity"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
)

func (p *Producer) CreateEvent(ctx context.Context, l entity.Link) error {
	ctx, span := tracer.Start(ctx, "kafka_producer UpdateEvent")
	defer span.End()

	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(l.Alias),
		Value: []byte(l.URL),
	})
	if err != nil {
		return fmt.Errorf("p.writer.WriteMessages: %w", err)
	}

	return nil
}
