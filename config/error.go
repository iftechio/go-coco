package config

import "github.com/pkg/errors"

var (
	ErrMissingEnvDefault = errors.New("missing default env config")
	ErrParseEnvFail      = errors.New("fail to parse env config")
)
