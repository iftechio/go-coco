package cmd

import (
	"os/exec"
	"path"

	"github.com/xieziyu/go-coco/cmd/tpl"
)

type SvcApps struct {
	Http    bool
	Grpc    bool
	Cronjob bool
	Looper  bool
}

type SvcInfras struct {
	Mongo  bool
	Redis  bool
	Sentry bool
}

// Service contains name and paths to services.
type Service struct {
	PkgName       string
	AbsolutePath  string
	AppName       string
	SvcPath       string
	GoVersion     string
	AlpineVersion string
	Description   string
	Apps          SvcApps
	Infras        SvcInfras
	WithInfra     bool
	WithProto     bool
	WithServer    bool
}

var apps = []string{
	"http-server",
	"grpc-server",
	"cronjob",
	"looper",
}

var infras = []string{
	"mongo",
	"redis",
	"sentry",
}

var protoMods = []string{
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest",
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest",
	"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
	"google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
}

var requiredMods = []string{
	"github.com/xieziyu/go-coco",
	"github.com/xieziyu/go-coco/utils/logger",
	"github.com/google/wire",
	"go.uber.org/automaxprocs",
}

func (s *Service) Create() error {
	internalPath := path.Join(s.AbsolutePath, "internal")
	configPath := path.Join(internalPath, "config")
	if err := ensureDir(s.AbsolutePath); err != nil {
		return err
	}
	if err := ensureDir(internalPath); err != nil {
		return err
	}
	if err := ensureDir(configPath); err != nil {
		return err
	}
	if err := createFile(s.AbsolutePath, "Dockerfile", tpl.Dockerfile(), s); err != nil {
		return err
	}
	if err := createFile(s.AbsolutePath, "main.go", tpl.MainGo(), s); err != nil {
		return err
	}
	if err := createFile(s.AbsolutePath, "Makefile", tpl.Makefile(), s); err != nil {
		return err
	}
	if err := createFile(s.AbsolutePath, "README.md", tpl.Readme(), s); err != nil {
		return err
	}
	if err := createFile(s.AbsolutePath, "wire.go", tpl.WireGo(), s); err != nil {
		return err
	}
	if err := createFile(configPath, "config.go", tpl.ConfigGo(), s); err != nil {
		return err
	}
	if err := createFile(configPath, "config.toml", tpl.ConfigToml(), s); err != nil {
		return err
	}
	if s.WithInfra {
		infraPath := path.Join(internalPath, "infra")
		if err := ensureDir(infraPath); err != nil {
			return err
		}
		if s.Infras.Mongo {
			if err := createFile(infraPath, "mongo.go", tpl.InfraMongo(), s); err != nil {
				return err
			}
		}
		if s.Infras.Redis {
			if err := createFile(infraPath, "redis.go", tpl.InfraRedis(), s); err != nil {
				return err
			}
		}
		if s.Infras.Sentry {
			if err := createFile(infraPath, "sentry.go", tpl.InfraSentry(), s); err != nil {
				return err
			}
		}
	}
	if s.Apps.Http {
		p := path.Join(internalPath, "http")
		if err := ensureDir(p); err != nil {
			return err
		}
		if err := createFile(p, "http.go", tpl.AppHttp(), s); err != nil {
			return err
		}
	}
	if s.Apps.Grpc {
		grpcPath := path.Join(internalPath, "grpc")
		if err := ensureDir(grpcPath); err != nil {
			return err
		}
		if err := createFile(grpcPath, "grpc.go", tpl.AppGrpc(), s); err != nil {
			return err
		}
		if s.WithProto {
			p := path.Join(grpcPath, "public")
			if err := ensureDir(p); err != nil {
				return err
			}
			if err := createFile(p, "public.go", tpl.GrpcPublicServer(), s); err != nil {
				return err
			}
		}
	}
	if s.Apps.Cronjob {
		p := path.Join(internalPath, "cronjob")
		if err := ensureDir(p); err != nil {
			return err
		}
		if err := createFile(p, "cronjob.go", tpl.AppCronjob(), s); err != nil {
			return err
		}
	}
	if s.Apps.Looper {
		p := path.Join(internalPath, "looper")
		if err := ensureDir(p); err != nil {
			return err
		}
		if err := createFile(p, "looper.go", tpl.AppLooper(), s); err != nil {
			return err
		}
	}
	if s.WithProto {
		apiPath := path.Join(s.AbsolutePath, "api")
		publicPath := path.Join(apiPath, "public")
		if err := ensureDir(apiPath); err != nil {
			return err
		}
		if err := ensureDir(publicPath); err != nil {
			return err
		}
		if err := createFile(apiPath, "Makefile", tpl.ProtoMakefile(), s); err != nil {
			return err
		}
		if err := createFile(apiPath, "buf.yaml", tpl.ProtoBufYaml(), s); err != nil {
			return err
		}
		if err := createFile(apiPath, "buf.gen.yaml", tpl.ProtoBufGenYaml(), s); err != nil {
			return err
		}
		if err := createFile(publicPath, "public.proto", tpl.PublicProto(), s); err != nil {
			return err
		}
		if err := bufModupdate(apiPath); err != nil {
			return err
		}
		if err := makeProto(s.AbsolutePath); err != nil {
			return err
		}
	}
	// TODO: Add more files
	return nil
}

func bufModupdate(path string) error {
	report(flagSuccess, "updating buf mod...")
	updateCmd := exec.Command("buf", "mod", "update")
	updateCmd.Dir = path
	if err := updateCmd.Run(); err != nil {
		return err
	}
	return nil
}

func makeProto(path string) error {
	report(flagSuccess, "generating proto files...")
	makeCmd := exec.Command("make", "proto")
	makeCmd.Dir = path
	if err := makeCmd.Run(); err != nil {
		return err
	}
	return nil
}

func makeWire(path string) error {
	report(flagSuccess, "making wire...")
	wireCmd := exec.Command("make", "wire")
	wireCmd.Dir = path
	if err := wireCmd.Run(); err != nil {
		return err
	}
	return nil
}
