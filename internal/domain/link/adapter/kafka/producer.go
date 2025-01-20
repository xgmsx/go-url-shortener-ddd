package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/entity"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
)

type Producer struct {
	writer *kafka.Writer
}

func New(writer *kafka.Writer) *Producer {
	return &Producer{writer: writer}
}

func (p *Producer) SendLink(ctx context.Context, l entity.Link) error {
	ctx, span := tracer.Start(ctx, "kafka SendLink")
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
