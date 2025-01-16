package kafka_consumer //nolint:stylecheck

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/dto"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/usecase"
	kafkaReader "github.com/xgmsx/go-url-shortener-ddd/pkg/kafka/reader"
)

type Consumer struct {
	uc     *usecase.UseCase
	stop   context.CancelFunc
	notify chan error
}

func New(ch chan error, reader *kafkaReader.Reader, uc *usecase.UseCase) *Consumer {
	ctx, stop := context.WithCancel(context.Background())
	c := &Consumer{
		uc:     uc,
		stop:   stop,
		notify: ch,
	}

	go func() {
	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			default:
				m, err := reader.FetchMessage(ctx)
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

				output, err := uc.CreateLink(ctx, input)
				if err != nil {
					log.Error().Err(err).Msg("kafka_consumer: uc.CreateLink")
					continue
				}
				log.Info().Msg("Link created: " + output.Str())

				if err = reader.CommitMessages(ctx, m); err != nil {
					log.Error().Err(err).Msg("kafka_consumer: reader.CommitMessages")
					continue
				}
			}
		}
	}()

	log.Info().Msg("Kafka consumer started")

	return c
}

func (c *Consumer) Close() {
	c.stop()
	log.Info().Msg("Kafka consumer closed")
}

func (c *Consumer) Notify(err error) {
	if err != nil {
		c.notify <- err
	}
}
