package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/dto"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/entity"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
)

func (u *UseCase) GetLink(ctx context.Context, input dto.GetLinkInput) (dto.GetLinkOutput, error) {
	ctx, span := tracer.Start(ctx, "usecase GetLink")
	defer span.End()

	var output dto.GetLinkOutput

	link, err := u.cache.GetLink(ctx, input.Alias)
	switch {
	case err == nil:
		return output.Load(link), nil
	case !errors.Is(err, entity.ErrNotFound):
		log.Error().Err(err).Msg("u.cache.GetLink")
	}

	link, err = u.db.GetLink(ctx, input.Alias, "")
	if err != nil {
		return output, fmt.Errorf("u.db.GetLink: %w", err)
	}

	err = u.cache.PutLink(ctx, link)
	if err != nil {
		log.Error().Err(err).Msg("u.cache.PutLink")
	}

	return output.Load(link), nil
}
