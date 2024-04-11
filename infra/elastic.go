package infra

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
)

type Elastic struct {
	*elasticsearch.Client
	Coco
}

type ElasticConfig struct {
	Addresses []string
	Username  string
	Password  string
}

// NewElastic provides a new Elastic client
func NewElastic(cfg ElasticConfig) (*Elastic, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err := es.Info()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	return &Elastic{Client: es}, nil
}
