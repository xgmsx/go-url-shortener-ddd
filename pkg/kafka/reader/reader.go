package reader

import (
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	Addr  []string `env:"KAFKA_BROKERS, required"`
	Topic string   `env:"KAFKA_INPUT_TOPIC, required"`
	Group string   `env:"KAFKA_GROUP, required"`
}

type Reader struct {
	*kafka.Reader
}

func New(c Config) (*Reader, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.Addr,
		GroupID:  c.Group,
		Topic:    c.Topic,
		MaxBytes: 10e6, // 10MB
	})

	return &Reader{Reader: r}, nil
}

func (r *Reader) Close() {
	err := r.Reader.Close()
	if err != nil {
		log.Error().Err(err).Msg("kafka - r.Reader.Close")
	}

	log.Info().Msg("Kafka reader closed")
}
