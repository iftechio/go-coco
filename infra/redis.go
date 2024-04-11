package infra

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmgoredisv8/v2"
)

type Redis struct {
	*redis.Client
	Coco
}

type RedisConfig struct {
	URL            string
	ConnectTimeout time.Duration // 初始化连接超时时间
}

// NewRedis provides a new Redis client
func NewRedis(cfg RedisConfig) (*Redis, error) {
	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	client.AddHook(apmgoredisv8.NewHook())

	ctx := context.Background()
	if cfg.ConnectTimeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, cfg.ConnectTimeout)
		defer cancel()
	}
	if err = client.Ping(ctx).Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return &Redis{Client: client}, nil
}
