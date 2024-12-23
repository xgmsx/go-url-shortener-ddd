package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"url-shortener/internal/config"
	kafkaReader "url-shortener/pkg/kafka/reader"
	kafkaWriter "url-shortener/pkg/kafka/writer"
	"url-shortener/pkg/postgres"
	"url-shortener/pkg/redis"

	"github.com/rs/zerolog/log"
)

type Dependencies struct {
	Redis       *redis.Client
	Postgres    *postgres.Pool
	KafkaWriter *kafkaWriter.Writer
	KafkaReader *kafkaReader.Reader
}

func Run(ctx context.Context, c config.Config) error {
	var (
		deps Dependencies
		err  error
	)

	// init dependencies
	deps.Postgres, err = postgres.New(ctx, c.Postgres)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}
	defer deps.Postgres.Close()

	deps.Redis, err = redis.New(c.Redis)
	if err != nil {
		return fmt.Errorf("redis.New: %w", err)
	}
	defer deps.Redis.Close()

	deps.KafkaWriter, err = kafkaWriter.New(c.KafkaWriter)
	if err != nil {
		return fmt.Errorf("kafkaWriter.New: %w", err)
	}
	defer deps.KafkaWriter.Close()

	deps.KafkaReader, err = kafkaReader.New(c.KafkaReader)
	if err != nil {
		return fmt.Errorf("kafkaReader.New: %w", err)
	}
	defer deps.KafkaReader.Close()

	// init domain
	ch := make(chan error, 4)
	defer close(ch)

	uc := getUCLink(deps)

	http := getHTTPController(ch, c, uc)
	defer http.Close()

	grpc := getGRPCController(ch, c, uc)
	defer grpc.Close()

	kafka := getKafkaController(ch, deps.KafkaReader, uc)
	defer kafka.Close()

	return waiting(ch)
}

func waiting(ch chan error) error {
	log.Info().Msg("App started")
	defer log.Info().Msg("App stopping...")

	ctxTerm, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	select {
	case <-ctxTerm.Done():
		log.Info().Msg("App got termination signal")
		return nil
	case err := <-ch:
		log.Info().Err(err).Msg("App got notify")
		return err
	}
}
