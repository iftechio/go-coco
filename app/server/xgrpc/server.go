package xgrpc

import (
	"net"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	*grpc.Server
}

type ServerConfig struct{}

type ServerConfigBuilder func(*ServerConfig)

// New provides a new gRPC server
func New(configs ...ServerConfigBuilder) *Server {
	cfg := ServerConfig{}
	for _, config := range configs {
		config(&cfg)
	}

	interceptors := []grpc.UnaryServerInterceptor{
		apmgrpc.NewUnaryServerInterceptor(),
		grpcrecovery.UnaryServerInterceptor(),
		sentryUnaryInterceptor(),
	}

	srv := &Server{
		Server: grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...)),
	}
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	return srv
}

// Start the gRPC server
func (s *Server) Start(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return s.Serve(l)
}
