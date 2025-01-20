package fetch

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/dto"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/entity"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
)

type Usecase struct {
	database database
	cache    cache
}

func New(d database, c cache) Usecase {
	return Usecase{database: d, cache: c}
}

func (u *Usecase) Fetch(ctx context.Context, input dto.FetchLinkInput) (dto.FetchLinkOutput, error) {
	ctx, span := tracer.Start(ctx, "usecase FetchLink")
	defer span.End()

	var output dto.FetchLinkOutput

	link, err := u.cache.GetLink(ctx, input.Alias)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		log.Error().Err(err).Msg("u.cache.GetLink")
	}

	if errors.Is(err, entity.ErrNotFound) {
		link, err = u.database.FindLink(ctx, input.Alias, "")
		if err != nil {
			return output, fmt.Errorf("u.database.FindLink: %w", err)
		}

		err = u.cache.PutLink(ctx, link)
		if err != nil {
			log.Error().Err(err).Msg("u.cache.PutLink")
		}

		return output.Load(link), nil
	}
	return output.Load(link), nil
}
