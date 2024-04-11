package infra

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/iftechio/go-coco/utils/logger"
	"github.com/pkg/errors"
)

type Postgres struct {
	*pg.DB
	Coco
}

type PostgresConfig struct {
	Addr     string
	User     string
	Password string
	Database string
}

// NewPostgres provides a new Postgres client
func NewPostgres(cfg PostgresConfig) (*Postgres, func(), error) {
	db := pg.Connect(&pg.Options{
		Addr:     cfg.Addr,
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Database,
	})
	cleanup := func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
	}

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		return nil, cleanup, errors.WithStack(err)
	}

	return &Postgres{DB: db}, cleanup, nil
}
