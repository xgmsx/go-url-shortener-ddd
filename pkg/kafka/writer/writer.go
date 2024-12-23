package writer

import (
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	Addr  []string `env:"KAFKA_BROKERS, required"`
	Topic string   `env:"KAFKA_OUTPUT_TOPIC, required"`
}

type Writer struct {
	*kafka.Writer
}

func New(c Config) (*Writer, error) {
	w := &kafka.Writer{
		Addr:     kafka.TCP(c.Addr...),
		Topic:    c.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Writer{Writer: w}, nil
}

func (w *Writer) Close() {
	err := w.Writer.Close()
	if err != nil {
		log.Error().Err(err).Msg("kafka - w.Writer.Close")
	}

	log.Info().Msg("Kafka writer closed")
}
