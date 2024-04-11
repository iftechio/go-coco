package looper

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// RedisLock implemented the DistributedLock interface
type RedisLock struct {
	redis *redis.Client
}

func NewRedisLock(client *redis.Client) *RedisLock {
	if client == nil {
		panic("NewRedisLock: redis client is nil")
	}
	return &RedisLock{
		redis: client,
	}
}

// Lock the key
func (l *RedisLock) Lock(ctx context.Context, key string, expiration time.Duration) (ok bool, err error) {
	now := time.Now().UnixMilli()
	ok, err = l.redis.SetNX(ctx, key, now, expiration).Result()
	if err != nil {
		return false, errors.WithStack(err)
	}
	return
}

// Unlock the key
func (l *RedisLock) Unlock(ctx context.Context, key string) (err error) {
	err = l.redis.Del(ctx, key).Err()
	return errors.WithStack(err)
}

// Renew the key expiration time
func (l *RedisLock) Renew(ctx context.Context, key string, expiration time.Duration) (err error) {
	err = l.redis.Expire(ctx, key, expiration).Err()
	return errors.WithStack(err)
}
