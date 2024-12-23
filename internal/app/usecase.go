package app

import (
	"url-shortener/internal/shortener/adapter/kafka_producer"
	"url-shortener/internal/shortener/adapter/postgres"
	"url-shortener/internal/shortener/adapter/redis"
	"url-shortener/internal/shortener/usecase"
)

func getUCLink(d Dependencies) *usecase.UseCase {
	return usecase.New(
		postgres.New(d.Postgres.Pool),
		redis.New(d.Redis.Client),
		kafka_producer.New(d.KafkaWriter.Writer),
	)
}
