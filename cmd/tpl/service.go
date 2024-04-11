package tpl

func Dockerfile() string {
	return `# Dockerfile contents...`
}

func MainGo() string {
	return `package main

import (
	_ "go.uber.org/automaxprocs"

	logger "github.com/iftechio/go-coco/utils/logger"
)

func main() {
	// 注入 apps
	apps, cleanup, err := injectAll()
	if err != nil {
		logger.Fatal("fatal error: inject apps. ", err)
		return
	}
	// 启动 apps
	if err = apps.Run(cleanup); err != nil {
		logger.Fatal("fatal error: apps run. ", err)
	}
}`
}

func Makefile() string {
	return `.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build: wire
	go build -o bin/server
	make tidy

{{- if .WithProto }}
.PHONY: proto
proto: 
	make -C api

.PHONY: all
all: proto build
{{ else }}
.PHONY: all
all: build
{{- end }}

.PHONY: test
test:
	NODE_ENV=test go test ./...

.PHONY: wire
wire:	
	wire

.PHONY: dev
dev: build
	NODE_ENV=dev bin/server
`
}

func Readme() string {
	return CodeBlockQuote + `# {{ .AppName | ToTitle }} Service

{{ .Description }}

## 项目结构

{{template "cbq"}}
{{ .AppName }}
├── internal
{{- if .Apps.Cronjob }}
│   ├── cronjob: 定时任务模块
{{- end }}
{{- if .Apps.Grpc }}
│   ├── grpc: GRPC Server
{{- end }}
{{- if .Apps.Http }}
│   ├── http: HTTP/HTTPS Server
{{- end }}
{{- if .WithProto }}
│   │   └── swagger: Protobuf 相关生成的 swagger 文件，请勿手动修改
{{- end }}
{{- if .WithInfra }}
│   ├── infra: 各类 Infra 组件
{{- end }}
{{- if .Apps.Looper }}
│   ├── looper: 循环任务模块 
{{- end }}
{{- if .WithProto }}
│   ├── proto: Protobuf 相关生成文件，请勿手动修改
{{- end }}
│   └── config: 服务配置
├── Dockerfile: 镜像构建配置
├── main.go
├── Makefile
├── README.md
├── wire_gen.go: wire 依赖注入生成文件，请勿手动修改
└── wire.go: wire 依赖注入定义
{{template "cbq"}}
`
}

func WireGo() string {
	return `//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"{{ .PkgName }}/internal/config"
{{- if .Apps.Cronjob }}
	"{{ .PkgName }}/internal/cronjob"
{{- end }}
{{- if .Apps.Grpc }}
	"{{ .PkgName }}/internal/grpc"
{{- end }}
{{- if .Apps.Http }}
	"{{ .PkgName }}/internal/http"
{{- end }}
{{- if .WithInfra }}
	"{{ .PkgName }}/internal/infra"
{{- end }}
{{- if .Apps.Looper }}
	"{{ .PkgName }}/internal/looper"
{{- end }}
	"github.com/iftechio/go-coco/app"
)

// AppManager is the manager of all the applications
type AppManager struct {
	app.Manager
	env string
}

// NewAppManager provides a new instance of Apps
func NewAppManager(
	cfg config.Config,
{{- if .Infras.Sentry }}
	sentry *infra.Sentry,
{{- end }}
{{- if .Infras.Redis }}
	redis *infra.Redis,
{{- end }}
{{- if .Infras.Mongo }}
	db *infra.MongoDB,
{{- end }}
{{- if .Apps.Http }}
	httpServer *http.Server,
{{- end }}
{{- if .Apps.Grpc }}
	grpcServer *grpc.Server,
{{- end }}
{{- if .Apps.Cronjob }}
	crons *cronjob.Jobs,
{{- end }}
{{- if .Apps.Looper }}
	loopers *looper.Jobs,
{{- end }}
) *AppManager {
	var mgr AppManager
	// NODE_ENV
	mgr.env = cfg.Env
	// infras
{{- if .Infras.Sentry }}
	mgr.RegisterInfra(sentry)
{{- end }}
{{- if .Infras.Redis }}
	mgr.RegisterInfra(redis)
{{- end }}
{{- if .Infras.Mongo }}
	mgr.RegisterInfra(db)
{{- end }}
	// apps
{{- if .Apps.Http }}
	mgr.RegisterApp(httpServer)
{{- end }}
{{- if .Apps.Grpc }}
	mgr.RegisterApp(grpcServer)
{{- end }}
{{- if .Apps.Cronjob }}
	mgr.RegisterApp(crons)
{{- end }}
{{- if .Apps.Looper }}
	mgr.RegisterApp(loopers)
{{- end }}
	return &mgr
}

// injectAll inject all the infras and apps
func injectAll() (*AppManager, func(), error) {
	panic(wire.Build(
		config.Provide,
{{- if .Infras.Sentry }}
		// sentry
		infra.NewSentry,
{{- end }}
{{- if .Infras.Redis }}
		// redis
		config.ProvideRedis,
		infra.NewRedis,
{{- end }}
{{- if .Infras.Mongo }}
		// mongo
		config.ProvideMongo,
		infra.NewMongoDB,
{{- end }}
		// apps
{{- if .Apps.Http }}
		http.NewServer,
{{- end }}
{{- if .Apps.Grpc }}
		grpc.NewServer,
{{- end }}
{{- if .Apps.Cronjob }}
		cronjob.NewJobs,
{{- end }}
{{- if .Apps.Looper }}
		looper.NewJobs,
{{- end }}
		NewAppManager,
		// Add service providers:
	))
}	
`
}

func ConfigGo() string {
	return BackquoteDef + `package config

import (
	_ "embed" // embed config toml

	"github.com/iftechio/go-coco/config"
)

//go:embed config.toml
var configToml string

type Config struct {
	Env                string                 {{template "bq"}}env:"NODE_ENV"{{template "bq"}}
{{- if .Apps.Cronjob }}
	Cronjob            string                 {{template "bq"}}env:"CRONJOB"{{template "bq"}}
{{- end }}
{{- if .Apps.Looper }}
	Looper             string                 {{template "bq"}}env:"LOOPER"{{template "bq"}}
{{- end }}
{{- if .WithServer }}
	AccessLog          string                 {{template "bq"}}env:"ACCESS_LOG"{{template "bq"}}
{{- end }}
{{- if .Apps.Http }}
	HTTPAddr           string                 {{template "bq"}}toml:"httpAddr"{{template "bq"}}
	RecordRequestBody  string                 {{template "bq"}}env:"RECORD_REQUEST_BODY"{{template "bq"}}
{{- end }}
{{- if .Apps.Grpc }}
	GRPCAddr           string                 {{template "bq"}}toml:"grpcAddr"{{template "bq"}}
{{- end }}
{{- if .Infras.Sentry }}
	SentryDSN          string                 {{template "bq"}}toml:"sentryDsn"{{template "bq"}}
{{- end }}
{{- if .Infras.Mongo }}
	DefaultMongo       map[string]MongoConfig {{template "bq"}}toml:"mongo"{{template "bq"}}
{{- end }}
{{- if .Infras.Redis }}
	DefaultRedis       map[string]RedisConfig {{template "bq"}}toml:"redis"{{template "bq"}}
{{- end }}
}
{{ if .Infras.Mongo }}
// Mongo Config
type MongoConfig struct {
	URL  string {{template "bq"}}toml:"url" env:"MONGO_URL"{{template "bq"}} // TODO: Config Mongo URL Env
	DB   string {{template "bq"}}toml:"db"{{template "bq"}}
	MaxPoolSize uint64 {{template "bq"}}toml:"maxPoolSize" env:"MONGO_MAX_POOL_SIZE"{{template "bq"}}
}
{{ end }}
{{- if .Infras.Redis }}
// Redis Config
type RedisConfig struct {
	URL string {{template "bq"}}toml:"url" env:"REDIS_URL"{{template "bq"}} // TODO: Config Redis URL Env
}
{{ end }}
// Provide config
func Provide() Config {
	cfg, err := config.Parse[Config](configToml)
	if err != nil {
		panic(err)
	}
	return cfg
}
{{ if .Infras.Mongo }}
// Provide MongoConfig
func ProvideMongo(cfg Config) MongoConfig {
	conf, err := config.ParseEnv(cfg.Env, cfg.DefaultMongo)
	if err != nil {
		panic(err)
	}
	return conf
}
{{ end }}
{{- if .Infras.Redis }}
// Provide RedisConfig
func ProvideRedis(cfg Config) RedisConfig {
	conf, err := config.ParseEnv(cfg.Env, cfg.DefaultRedis)
	if err != nil {
		panic(err)
	}
	return conf
}
{{ end }}
`
}

func ConfigToml() string {
	return `{{- if .Apps.Http }}httpAddr = ":3000"{{- end }}
{{ if .Apps.Grpc }}grpcAddr = ":9090"{{- end }}
{{ if .Infras.Sentry }}sentryDsn = "" # TODO: Config Sentry DSN{{- end }}
{{ if .Infras.Mongo }}
[mongo]
[mongo.default]
url = "mongodb://localhost:27017"
db = "yourDB"
maxPoolSize = 50
[mongo.ci]
url = "mongodb://mongo:27017"
db = "yourDB-ci"
maxPoolSize = 50
[mongo.test]
url = "mongodb://localhost:27017"
db = "yourDB-test"
maxPoolSize = 50
{{- end }}
{{ if .Infras.Redis }}
[redis]
[redis.default]
url = "redis://localhost:6379"
[redis.ci]
url = "redis://redis:6379"
[redis.test]
url = "redis://localhost:6379"
{{- end }}
`
}
