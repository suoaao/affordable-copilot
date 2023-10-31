package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"time"
)

type CoTokenCache struct {
	cache     *cache.Cache
	keyFormat string
}

func NewCoTokenCache(rdb *redis.Client) *CoTokenCache {
	return &CoTokenCache{
		cache:     cache.New(&cache.Options{Redis: rdb}),
		keyFormat: "copilot_internal/v2/token/%s",
	}
}

func (r *CoTokenCache) Get(ctx context.Context, authToken string) ([]byte, error) {
	var coToken []byte
	key := fmt.Sprintf(r.keyFormat, authToken)
	err := r.cache.Get(ctx, key, &coToken)
	return coToken, err
}

func (r *CoTokenCache) Set(ctx context.Context, authToken string, coToken []byte, ttl time.Duration) error {
	if len(coToken) == 0 {
		return errors.New("empty coToken")
	}
	key := fmt.Sprintf(r.keyFormat, authToken)
	return r.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: coToken,
		TTL:   ttl,
	})
}
