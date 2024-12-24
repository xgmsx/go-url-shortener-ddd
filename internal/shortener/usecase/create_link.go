package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/dto"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/entity"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
)

func (u *UseCase) CreateLink(ctx context.Context, input dto.CreateLinkInput) (dto.CreateLinkOutput, error) {
	ctx, span := tracer.Start(ctx, "usecase CreateLink")
	defer span.End()

	var (
		output dto.CreateLinkOutput
		id     = uuid.New()
		alias  = base64.RawURLEncoding.EncodeToString(id[:])
	)

	link, err := u.db.GetLink(ctx, "", input.URL)
	if err == nil {
		return output.Load(link), entity.ErrAlreadyExist
	}

	link = entity.Link{
		ID:        id,
		URL:       input.URL,
		Alias:     alias,
		ExpiredAt: time.Now().Add(linkTTL),
	}

	err = u.db.CreateLink(ctx, link)
	if err != nil {
		return output, fmt.Errorf("u.db.CreateLink: %w", err)
	}

	err = u.cache.PutLink(ctx, link)
	if err != nil {
		log.Error().Err(err).Msg("u.cache.PutLink")
	}

	err = u.broker.CreateEvent(ctx, link)
	if err != nil {
		log.Error().Err(err).Msg("u.broker.CreateEvent")
	}

	return output.Load(link), nil
}
