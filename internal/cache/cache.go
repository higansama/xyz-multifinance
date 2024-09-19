package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/higansama/xyz-multi-finance/internal/redis"
	"github.com/pkg/errors"
)

type Store interface {
	Get(ctx context.Context, key string, value any) (bool, error)
	Put(ctx context.Context, key string, payload any, ttl time.Duration) error
	Increment(ctx context.Context, key string, value int64) (int64, error)
	Remove(ctx context.Context, key string) error
}

type Cache struct {
	store Store
}

func (c *Cache) Get(ctx context.Context, key string, value any) (bool, error) {
	return c.store.Get(ctx, key, value)
}

func (c *Cache) Put(ctx context.Context, key string, payload any, ttl time.Duration) error {
	return c.store.Put(ctx, key, payload, ttl)
}

func (c *Cache) Forever(ctx context.Context, key string, payload any) error {
	return c.store.Put(ctx, key, payload, time.Hour*1000*200)
}

func (c *Cache) Increment(ctx context.Context, key string, value int64) (int64, error) {
	return c.store.Increment(ctx, key, value)
}

func (c *Cache) Remove(ctx context.Context, key string) error {
	return c.store.Remove(ctx, key)
}

func NewCache(store Store) *Cache {
	return &Cache{store: store}
}

type RedisStore struct {
	prefix string
	client *cache.Cache
}

func (c *RedisStore) Get(ctx context.Context, key string, value any) (bool, error) {
	err := c.client.Get(ctx, c.prefix+key, value)
	if err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			return false, nil
		}
		return false, errors.WithStack(err)
	}

	return true, nil
}

func (c *RedisStore) Put(ctx context.Context, key string, payload any, ttl time.Duration) error {
	err := c.client.Set(&cache.Item{
		Ctx:   ctx,
		Key:   c.prefix + key,
		Value: payload,
		TTL:   ttl,
	})
	return errors.WithStack(err)
}

func (c *RedisStore) Increment(ctx context.Context, key string, value int64) (int64, error) {
	var v int64
	err := c.client.Get(ctx, c.prefix+key, &v)
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		return 0, errors.WithStack(err)
	}

	v = v + value
	err = c.client.Set(&cache.Item{
		Ctx:   ctx,
		Key:   c.prefix + key,
		Value: v,
	})
	if err != nil {
		return 0, errors.Wrap(err, "cache: failed when incrementing the value")
	}

	return v, nil
}

func (c *RedisStore) Remove(ctx context.Context, key string) error {
	err := c.client.Delete(ctx, c.prefix+key)
	return errors.WithStack(err)
}

func NewRedisStore(client *redis.Client) *RedisStore {
	cacheClient := cache.New(&cache.Options{
		Redis: client.Instance,
	})
	return &RedisStore{client: cacheClient, prefix: client.Prefix}
}
