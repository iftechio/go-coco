package tpl

func InfraSentry() string {
	return `package infra

import (
	"github.com/getsentry/sentry-go"
	"{{ .PkgName }}/internal/config"
	"github.com/xieziyu/go-coco/infra"
)

type Sentry = infra.Sentry

// NewSentry provides a new instance of Sentry
func NewSentry(cfg config.Config) (*Sentry, error) {
	return infra.NewSentry(sentry.ClientOptions{
		Dsn:         cfg.SentryDSN,
		Environment: cfg.Env,
	})
}
`
}

func InfraMongo() string {
	return `package infra

import (
	"{{ .PkgName }}/internal/config"
	"github.com/xieziyu/go-coco/infra"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	*infra.Mongo
	MyCol *mongo.Collection // TODO: Naming Your Collection
}

// NewMongoDB provides a new instance of mongo db client
func NewMongoDB(cfg config.MongoConfig) (*MongoDB, error) {
	client, err := infra.NewMongo(infra.MongoConfig{
		URL:         cfg.URL,
		MaxPoolSize: cfg.MaxPoolSize,
	})
	if err != nil {
		return nil, err
	}
	return &MongoDB{
		Mongo: client,
		MyCol: client.Database(cfg.DB).Collection("mycols"), // TODO: Naming Your Collection
	}, nil
}
`
}

func InfraRedis() string {
	return `package infra

import (
	"{{ .PkgName }}/internal/config"
	"github.com/xieziyu/go-coco/infra"
)

type Redis struct {
	*infra.Redis
}

// NewRedis provides a new instance of redis
func NewRedis(cfg config.RedisConfig) (*Redis, error) {
	client, err := infra.NewRedis(infra.RedisConfig{
		URL: cfg.URL,
	})
	if err != nil {
		return nil, err
	}
	return &Redis{
		Redis:      client,
	}, nil
}
`
}
