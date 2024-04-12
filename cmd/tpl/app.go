package tpl

func AppHttp() string {
	return `package http

import (
	"context"
{{- if .WithProto }}
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
{{- end }}

	"{{ .PkgName }}/internal/config"
{{- if .WithProto }}
	"{{ .PkgName }}/internal/proto/public"
{{- end }}
	_ "github.com/iftechio/go-coco/app/server/profile"
	"github.com/iftechio/go-coco/app/server/xecho"
	"github.com/iftechio/go-coco/app/server/xecho/middleware"
{{- if .WithProto }}
	"github.com/iftechio/go-coco/app/server/xgrpc"
{{- end }}
	logger "github.com/iftechio/go-coco/utils/logger"
)

type Server struct {
	srv  *xecho.Server
	conf config.Config
}

{{- if .WithProto }}
//go:embed swagger
var swagger embed.FS
{{- end }}

// NewServer provides a new http server
func NewServer(conf config.Config) (*Server, func(), error) {
	srv := xecho.New()
	return &Server{
			srv:  srv,
			conf: conf,
		}, func() {
			if err := srv.Shutdown(context.Background()); err != nil {
				logger.Error(err)
			}
		}, nil
}

// Start the http server
func (s *Server) Start() error {
	s.srv.Use(middleware.HealthProbe())
{{- if .Infras.Sentry }}
	s.srv.Use(middleware.Sentry())
{{- end }}
	// TODO: Add Routings
{{- if .WithProto }}
	s.srv.Use(xgrpc.GatewayMiddleware(
		public.RegisterPublicServiceHandlerFromEndpoint,
		s.conf.GRPCAddr,
		"/1.0",
	))
{{ end }}
{{- if .WithProto }}
	// Swagger Doc
	s.srv.GET("/swagger/*", echo.WrapHandler(http.FileServer(http.FS(swagger))))
{{- end }}
	return s.srv.Start(s.conf.HTTPAddr)
}

// IsEnabled checks if the server enabled
func (s *Server) IsEnabled() bool {
{{- if .Apps.Cronjob }}
	// cronjob 启动时默认不启动 server
	return s.conf.Cronjob == ""
{{- else }}
	return true
{{- end }}
}
`
}

func GrpcPublicServer() string {
	return `package public

import (
	"context"

	"google.golang.org/grpc"

	"{{ .PkgName }}/internal/proto/public"
)

type Server struct {
	public.UnsafePublicServiceServer
}

// New provides public service server
func New() *Server {
	return &Server{}
}

// Register public service server
func (s *Server) Register(srv grpc.ServiceRegistrar) {
	public.RegisterPublicServiceServer(srv, s)
}

func (s *Server) Echo(ctx context.Context, req *public.EchoRequest) (*public.EchoResponse, error) {
	return &public.EchoResponse{Data: req.GetData()}, nil
}
`
}

func AppGrpc() string {
	return `package grpc

import (
	"{{ .PkgName }}/internal/config"
	"github.com/iftechio/go-coco/app/server/xgrpc"
{{- if .WithProto }}
	"{{ .PkgName }}/internal/grpc/public"
{{- end }}
)

type Server struct {
	srv  *xgrpc.Server
	conf config.Config
}

// NewServer provides a new gRPC server
func NewServer(
	conf config.Config,
) (*Server, func(), error) {
	srv := xgrpc.New()

	// TODO: Register Handlers
{{- if .WithProto }}
	public.New().Register(srv)
{{- end }}

	return &Server{
			srv:  srv,
			conf: conf,
		}, func() {
			srv.GracefulStop()
		}, nil
}

// Start the gRPC server
func (s *Server) Start() error {
	return s.srv.Start(s.conf.GRPCAddr)
}

// IsEnabled checks if the server enabled
func (s *Server) IsEnabled() bool {
{{- if .Apps.Cronjob }}
	// cronjob 启动时默认不启动 server
	return s.conf.Cronjob == ""
{{- else }}
	return true
{{- end }}
}	
`
}

func AppCronjob() string {
	return `package cronjob

import (
	"context"
{{ if .Infras.Sentry }}
	"github.com/getsentry/sentry-go"
{{- end }}
	"{{ .PkgName }}/internal/config"
	logger "github.com/iftechio/go-coco/utils/logger"
)

type Jobs struct {
	conf config.Config
}

func NewJobs(conf config.Config) *Jobs {
	return &Jobs{
		conf: conf,
	}
}

// Start a cronjob
func (j *Jobs) Start() error {
	ctx := context.Background()
	log := logger.F(ctx)
	log.Infof("[%s] cronjob started", j.conf.Cronjob)
	var err error
	switch j.conf.Cronjob {
	// TODO: add cronjob routines:
	// case "MyCronJob":
	// 	 err = j.RunMyCronJob(context.Background())
	}
	if err != nil {
{{- if .Infras.Sentry }}
		sentry.CaptureException(err)
{{- end }}
		log.Error(err)
	}
	log.Infof("[%s] cronjob done.", j.conf.Cronjob)
	return nil
}

func (j *Jobs) IsEnabled() bool {
	return j.conf.Cronjob != ""
}	
`
}

func AppLooper() string {
	return `package looper

import (
	"context"
	"time"

	"{{ .PkgName }}/internal/config"
	"{{ .PkgName }}/internal/infra"
	"github.com/iftechio/go-coco/app/looper"
	"github.com/iftechio/go-coco/utils/sync/errgroup"
	logger "github.com/iftechio/go-coco/utils/logger"
)

type Jobs struct {
	conf  config.Config
	// TODO: Define Your Jobs
	myJob1 *looper.SinglePod
	myJob2 *looper.SinglePod
}

func NewJobs(
	conf config.Config,
	rds *infra.Redis,
) (w *Jobs, cleanup func(), err error) {
	// TODO: Define Your Jobs
	myJob1, cleanup1, err := looper.NewSinglePod(
		"{{ .AppName }}:my-job1",
		looper.NewRedisLock(rds.Client),
		func(ctx context.Context) error {
			// TODO: handler
			return nil
		},
		nil,
		looper.SinglePodConfig{
			Sleep:        time.Second,
			MaxRunTime:   time.Minute * 5,
			StandbySleep: time.Second * 10,
		})
	if err != nil {
		return nil, nil, err
	}
	myJob2, cleanup2, err := looper.NewSinglePod(
		"{{ .AppName }}:my-job2",
		looper.NewRedisLock(rds.Client),
		func(ctx context.Context) error {
			// TODO: handler
			return nil
		},
		nil,
		looper.SinglePodConfig{
			Sleep:        time.Second,
			MaxRunTime:   time.Minute * 5,
			StandbySleep: time.Second * 10,
		})
	w = &Jobs{
		conf:   conf,
		myJob1: myJob1,
		myJob2: myJob2,
	}
	cleanup = func() {
		cleanup1()
		cleanup2()
	}
	return
}

// Start jobs
func (j *Jobs) Start() (err error) {
	eg := errgroup.New()
	eg.Go(func() error {
		logger.Info("MyJob1 looper start")
		return j.myJob1.Start()
	})
	eg.Go(func() error {
		logger.Info("MyJob2 looper start")
		return j.myJob2.Start()
	})
	return eg.Wait()
}

func (c *Jobs) IsEnabled() bool {
	return c.conf.Looper == "ON"
}
`
}
