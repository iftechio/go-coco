package infra

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"

	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmmongo/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	*mongo.Client
	Coco
}

type MongoConfig struct {
	URL             string
	ConnectTimeout  time.Duration // 初始化连接超时时间
	SocketTimeout   time.Duration
	MaxConnIdleTime time.Duration
	MaxPoolSize     uint64
}

// NewMongo provides a new Mongo client
func NewMongo(c MongoConfig) (*Mongo, error) {
	opt := options.Client().ApplyURI(c.URL)

	if c.SocketTimeout > 0 {
		opt.SetSocketTimeout(c.SocketTimeout)
	} else {
		opt.SetSocketTimeout(5 * time.Second)
	}

	if c.MaxPoolSize > 0 {
		opt.SetMaxPoolSize(c.MaxPoolSize)
	} else {
		opt.SetMaxPoolSize(100)
	}

	if c.MaxConnIdleTime > 0 {
		opt.SetMaxConnIdleTime(c.MaxConnIdleTime)
	} else {
		opt.SetMaxConnIdleTime(30 * time.Second)
	}
	opt.SetMonitor(apmmongo.CommandMonitor())
	client, err := mongo.NewClient(opt)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	if c.ConnectTimeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, c.ConnectTimeout)
		defer cancel()
	}
	if err := client.Connect(ctx); err != nil {
		return nil, errors.WithStack(err)
	}
	return &Mongo{
		Client: client,
	}, client.Ping(ctx, readpref.Primary())
}

type TransactionFn = func(sessionContext mongo.SessionContext) (interface{}, error)

func (m *Mongo) WithTransaction(ctx context.Context, fn TransactionFn) (interface{}, error) {
	txnOpts := options.Transaction().
		SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
		SetReadConcern(readconcern.Snapshot()).
		SetReadPreference(readpref.Primary())
	session, err := m.Client.StartSession()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer session.EndSession(ctx)
	res, err := session.WithTransaction(ctx, fn, txnOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}
