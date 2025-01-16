package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract.go

type kafkaReader interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
}
