package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/entity"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
)

const ttl = time.Hour

type Redis struct {
	client *redis.Client
}

func New(client *redis.Client) *Redis {
	return &Redis{client: client}
}

func (r *Redis) PutLink(ctx context.Context, link entity.Link) error {
	ctx, span := tracer.Start(ctx, "redis PutLink")
	defer span.End()

	data, err := json.Marshal(link)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	err = r.client.Set(ctx, link.Alias, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("r.client.Set: %w", err)
	}

	return nil
}

func (r *Redis) GetLink(ctx context.Context, alias string) (entity.Link, error) {
	ctx, span := tracer.Start(ctx, "redis GetLink")
	defer span.End()

	var link entity.Link

	data, err := r.client.Get(ctx, alias).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return link, entity.ErrNotFound
		}

		return link, fmt.Errorf("r.client.Get: %w", err)
	}

	err = json.Unmarshal(data, &link)
	if err != nil {
		return link, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return link, nil
}
