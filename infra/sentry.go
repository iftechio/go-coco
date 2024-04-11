package infra

import "github.com/getsentry/sentry-go"

type Sentry struct {
	Coco
}

func NewSentry(opts sentry.ClientOptions) (*Sentry, error) {
	if err := sentry.Init(opts); err != nil {
		return nil, err
	}
	return &Sentry{}, nil
}
