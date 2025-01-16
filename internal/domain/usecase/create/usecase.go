package create

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/dto"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/entity"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
)

const linkTTL = 24 * time.Hour

type Usecase struct {
	database  database
	cache     cache
	publisher publisher
}

func New(d database, c cache, p publisher) Usecase {
	return Usecase{database: d, cache: c, publisher: p}
}

func (u *Usecase) Create(ctx context.Context, input dto.CreateLinkInput) (dto.CreateLinkOutput, error) {
	ctx, span := tracer.Start(ctx, "usecase CreateLink")
	defer span.End()

	var (
		id     = uuid.New()
		alias  = base64.RawURLEncoding.EncodeToString(id[:])
		output dto.CreateLinkOutput
	)

	link := entity.Link{
		ID:        id,
		URL:       input.URL,
		Alias:     alias,
		ExpiredAt: time.Now().Add(linkTTL),
	}

	err := u.database.CreateLink(ctx, link)
	if err != nil {
		return output, fmt.Errorf("u.database.CreateLink: %w", err)
	}

	err = u.cache.PutLink(ctx, link)
	if err != nil {
		return output, fmt.Errorf("u.cache.PutLink: %w", err)
	}

	err = u.publisher.SendLink(ctx, link)
	if err != nil {
		return output, fmt.Errorf("u.broker.SendLink: %w", err)
	}

	return output.Load(link), nil
}
