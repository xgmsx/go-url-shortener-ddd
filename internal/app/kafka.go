package app

import (
	"url-shortener/internal/shortener/controller/kafka_consumer"
	"url-shortener/internal/shortener/usecase"
	"url-shortener/pkg/kafka/reader"
)

func getKafkaController(ch chan error, reader *reader.Reader, uc *usecase.UseCase) *kafka_consumer.Consumer {
	return kafka_consumer.New(ch, reader, uc)
}
