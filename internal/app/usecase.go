package app

import (
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/adapter/kafka_producer"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/adapter/postgres"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/adapter/redis"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/usecase"
)

func getUCLink(d Dependencies) *usecase.UseCase {
	return usecase.New(
		postgres.New(d.Postgres.Pool),
		redis.New(d.Redis.Client),
		kafka_producer.New(d.KafkaWriter.Writer),
	)
}
