package app

import (
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/controller/kafka_consumer"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/usecase"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/kafka/reader"
)

func getKafkaController(ch chan error, r *reader.Reader, uc *usecase.UseCase) *kafka_consumer.Consumer {
	return kafka_consumer.New(ch, r, uc)
}
