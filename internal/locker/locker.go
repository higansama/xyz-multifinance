package locker

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/higansama/xyz-multi-finance/internal/redis"
	"github.com/pkg/errors"
)

var (
	ErrLockNotObtained = errors.New("lock not obtained")
	ErrLockNotHeld     = errors.New("lock not held")
)

type Client interface {
	Obtain(ctx context.Context, key string, ttl time.Duration) (Lock, error)
}

type Locker struct {
	client Client
}

func (l *Locker) Obtain(ctx context.Context, key string, ttl time.Duration, retryTtl time.Duration) (Lock, error) {
	// to make sure don't retry forever
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(ttl))
		defer cancel()
	}

	var ticker *time.Ticker
	for {
		lock, err := l.client.Obtain(ctx, key, ttl)
		if err != nil && !errors.Is(err, ErrLockNotObtained) {
			return nil, err
		} else if err == nil {
			return lock, nil
		}

		backoff := retryTtl
		if backoff < 1 {
			return nil, errors.Wrap(ErrLockNotObtained, "locker")
		}

		if ticker == nil {
			ticker = time.NewTicker(backoff)
			defer ticker.Stop()
		} else {
			ticker.Reset(backoff)
		}

		select {
		case <-ctx.Done():
			return nil, errors.Wrap(ctx.Err(), "locker")
		case <-ticker.C:
		}
	}
}

func NewLocker(client Client) *Locker {
	return &Locker{client: client}
}

type RedisLocker struct {
	prefix string
	client *redislock.Client
}

func (rl *RedisLocker) Obtain(ctx context.Context, key string, ttl time.Duration) (Lock, error) {
	lockClient, err := rl.client.Obtain(ctx, rl.prefix+key, ttl, nil)
	if err != nil {
		return nil, errors.Wrap(transformError(err), "locker")
	}

	return NewRedisLock(ctx, lockClient), nil
}

func NewRedisLocker(client *redis.Client) *RedisLocker {
	lockClient := redislock.New(client.Instance)
	return &RedisLocker{prefix: client.Prefix + "locker:", client: lockClient}
}

type Lock interface {
	Release() error
	Refresh(ttl time.Duration) error
}

type RedisLock struct {
	ctx    context.Context
	client *redislock.Lock
}

func (rl *RedisLock) Release() error {
	err := rl.client.Release(rl.ctx)
	return errors.Wrap(transformError(err), "locker")
}

func (rl *RedisLock) Refresh(ttl time.Duration) error {
	err := rl.client.Refresh(rl.ctx, ttl, nil)
	return errors.Wrap(transformError(err), "locker")
}

func NewRedisLock(ctx context.Context, client *redislock.Lock) *RedisLock {
	return &RedisLock{ctx: ctx, client: client}
}

func transformError(err error) error {
	if errors.Is(err, redislock.ErrNotObtained) {
		return ErrLockNotObtained
	} else if errors.Is(err, redislock.ErrLockNotHeld) {
		return ErrLockNotHeld
	}
	return err
}
