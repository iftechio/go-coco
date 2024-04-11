package config

import (
	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

// Parse 从toml和env中解析配置出内容
func Parse[Config any](content string) (Config, error) {
	var cfg Config
	if _, err := toml.Decode(content, &cfg); err != nil {
		return cfg, errors.WithStack(err)
	}
	if err := env.Parse(&cfg); err != nil {
		return cfg, errors.WithStack(err)
	}
	return cfg, nil
}

// ParseEnv 根据当前env解析配置内容，解析优先级如下
//  - toml中定义的对应env下的配置
//  - env中定义的环境变量配置
//  - toml中定义的default配置
func ParseEnv[Config any](currentEnv string, defaultConfig map[string]Config) (Config, error) {
	conf, ok := defaultConfig[currentEnv]
	if ok {
		return conf, nil
	}
	conf, ok = defaultConfig["default"]
	if !ok {
		return conf, errors.WithStack(ErrMissingEnvDefault)
	}
	if err := env.Parse(&conf); err != nil {
		return conf, errors.WithStack(ErrParseEnvFail)
	}
	return conf, nil
}
