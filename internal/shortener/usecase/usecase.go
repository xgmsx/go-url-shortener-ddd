package usecase

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/entity"
)

const linkTTL = 24 * time.Hour

type Database interface {
	CreateLink(ctx context.Context, tx pgx.Tx, l entity.Link) error
	GetLink(ctx context.Context, alias, url string) (entity.Link, error)
	Tx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

type Cache interface {
	PutLink(ctx context.Context, l entity.Link) error
	GetLink(ctx context.Context, alias string) (entity.Link, error)
}

type Broker interface {
	CreateEvent(ctx context.Context, l entity.Link) error
}

type UseCase struct {
	db     Database
	cache  Cache
	broker Broker
}

func New(db Database, cache Cache, broker Broker) *UseCase {
	return &UseCase{
		db:     db,
		cache:  cache,
		broker: broker,
	}
}
