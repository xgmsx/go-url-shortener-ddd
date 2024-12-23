package kafka_producer

import (
	"context"
	"fmt"

	"url-shortener/internal/shortener/entity"
	"url-shortener/pkg/observability/otel/tracer"

	"github.com/segmentio/kafka-go"
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
