package kafka

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/dto"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/usecase/create"
	kafkaReader "github.com/xgmsx/go-url-shortener-ddd/pkg/kafka/reader"
)

type Consumer struct {
	reader *kafkaReader.Reader
	uc     create.Usecase
}

func New(reader *kafkaReader.Reader, uc create.Usecase) *Consumer {
	return &Consumer{reader: reader, uc: uc}
}

func (c *Consumer) Consume(ctx context.Context) {
	log.Info().Msg("Kafka consumer started")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			m, err := c.reader.FetchMessage(ctx)
			if err != nil {
				log.Error().Err(err).Msg("kafka_consumer: reader.FetchMessage")
				continue
			}

			var input dto.CreateLinkInput

			err = json.Unmarshal(m.Value, &input)
			if err != nil {
				log.Error().Err(err).Msg("kafka_consumer: json.Unmarshal")
				continue
			}

			output, err := c.uc.Create(ctx, input)
			if err != nil {
				log.Error().Err(err).Msg("kafka_consumer: uc.CreateLink")
				continue
			}
			log.Info().Msg("Link created: " + output.Str())

			if err = c.reader.CommitMessages(ctx, m); err != nil {
				log.Error().Err(err).Msg("kafka_consumer: reader.CommitMessages")
				continue
			}
		}
	}
}
