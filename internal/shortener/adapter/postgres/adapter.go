package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/entity"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *Postgres {
	return &Postgres{
		pool: p,
	}
}

func (p *Postgres) CreateLink(ctx context.Context, link entity.Link) error {
	ctx, span := tracer.Start(ctx, "postgres CreateLink")
	defer span.End()

	dataset := goqu.Insert("links").Rows(goqu.Record{
		"id":         link.ID,
		"url":        link.URL,
		"alias":      link.Alias,
		"updated_at": time.Now(),
		"expired_at": link.ExpiredAt,
	})

	sql, _, err := dataset.ToSQL()
	if err != nil {
		return fmt.Errorf("dataset.ToSQL: %w", err)
	}

	_, err = p.pool.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	return nil
}

func (p *Postgres) GetLink(ctx context.Context, alias, url string) (entity.Link, error) {
	ctx, span := tracer.Start(ctx, "postgres GetLink")
	defer span.End()

	var link entity.Link

	dataset := goqu.
		Select("id", "url", "alias", "expired_at").
		From("links")

	switch {
	case alias != "" && url != "":
		dataset = dataset.Where(goqu.C("alias").Eq(alias), goqu.C("url").Eq(url))
	case alias != "":
		dataset = dataset.Where(goqu.C("alias").Eq(alias))
	case url != "":
		dataset = dataset.Where(goqu.C("url").Eq(url))
	default:
		return link, fmt.Errorf("query validation")
	}

	sql, _, err := dataset.ToSQL()
	if err != nil {
		return link, fmt.Errorf("dataset.ToSQL: %w", err)
	}

	row := p.pool.QueryRow(ctx, sql)
	if err := row.Scan(&link.ID, &link.URL, &link.Alias, &link.ExpiredAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return link, entity.ErrNotFound
		}
		return link, fmt.Errorf("row.Scan: %w", err)
	}

	return link, nil
}
