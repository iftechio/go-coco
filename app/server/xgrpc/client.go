package xgrpc

import (
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
)

func Connect(host string) (grpc.ClientConnInterface, error) {
	return grpc.Dial(
		host,
		grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(
			apmgrpc.NewUnaryClientInterceptor(),
		),
	)
}

func MustConnect(host string) grpc.ClientConnInterface {
	cc, err := Connect(host)
	if err != nil {
		panic(err)
	}
	return cc
}
