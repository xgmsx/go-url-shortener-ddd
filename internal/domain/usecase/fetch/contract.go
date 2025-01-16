package fetch

import (
	"context"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/entity"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract.go

type database interface {
	FindLink(ctx context.Context, alias string, url string) (*entity.Link, error)
}

type cache interface {
	GetLink(ctx context.Context, alias string) (*entity.Link, error)
	PutLink(context.Context, entity.Link) error
}
