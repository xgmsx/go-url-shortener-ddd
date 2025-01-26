package kafka

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/dto"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/create"
)

type Consumer struct {
	kafka kafkaReader
	uc    create.Usecase
}

func New(k kafkaReader, uc create.Usecase) *Consumer {
	return &Consumer{kafka: k, uc: uc}
}

func (c *Consumer) Consume(ctx context.Context) error {
	log.Info().Msg("Kafka consumer started")
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			m, err := c.kafka.FetchMessage(ctx)
			if err != nil {
				log.Error().Err(err).Msg("c.kafka.FetchMessage")
				continue
			}

			var input dto.CreateLinkInput
			err = json.Unmarshal(m.Value, &input)
			if err != nil {
				log.Error().Err(err).Msg("json.Unmarshal")
				continue
			}

			output, err := c.uc.Create(ctx, input)
			if err != nil {
				log.Error().Err(err).Msg("uc.CreateLink")
				continue
			}
			log.Info().Msg("Link created: " + output.Str())

			if err = c.kafka.CommitMessages(ctx, m); err != nil {
				log.Error().Err(err).Msg("c.kafka.CommitMessages")
				continue
			}
		}
	}
}
