package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Config struct {
	User     string `env:"POSTGRES_USER, required"`
	Password string `env:"POSTGRES_PASSWORD, required"`
	Port     string `env:"POSTGRES_PORT, default=5432"`
	Host     string `env:"POSTGRES_HOST, required"`
	DBName   string `env:"POSTGRES_DB, required"`
}

type Pool struct {
	*pgxpool.Pool
}

func New(ctx context.Context, c Config) (*Pool, error) {
	dsn := fmt.Sprintf("user=%s password=%s port=%s host=%s dbname=%s",
		c.User, c.Password, c.Port, c.Host, c.DBName)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	return &Pool{Pool: pool}, nil
}

func (p *Pool) Close() {
	p.Pool.Close()
	log.Info().Msg("Postgres closed")
}
