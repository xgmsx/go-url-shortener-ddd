package create

import (
	"context"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/entity"
)

type database interface {
	CreateLink(context.Context, entity.Link) error
}

type cache interface {
	PutLink(context.Context, entity.Link) error
}

type publisher interface {
	SendLink(ctx context.Context, link entity.Link) error
}
