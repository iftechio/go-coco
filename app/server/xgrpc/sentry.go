package xgrpc

import (
	"context"

	"github.com/getsentry/sentry-go"
	xLogger "github.com/iftechio/go-coco/utils/logger"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/iftechio/go-coco/app/server/custom"
)

func asCustomError(err error) (*custom.Error, bool) {
	customErr := &custom.Error{}
	if ok := errors.As(err, &customErr); ok {
		return customErr, true
	}
	return nil, false
}

func sentryUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		hub := sentry.CurrentHub().Clone()
		scope := hub.Scope()
		if info != nil {
			scope.SetExtra("method", info.FullMethod)
		}
		if req != nil {
			scope.SetExtra("request", req)
		}
		ctx = sentry.SetHubOnContext(ctx, hub)
		defer func() {
			if v := recover(); v != nil {
				hub.Recover(v)
				e, ok := v.(error)
				if !ok {
					e = errors.Errorf("%#v", e)
				}
				err = e // 保证 err 被带出去
				return
			}
			if err == nil {
				return
			}
			xLogger.F(ctx).Error(err)
			if e, ok := asCustomError(err); ok {
				if e.Err != nil {
					hub.CaptureException(e.Err)
				}
				return
			}
			switch status.Code(err) {
			case codes.Unknown, codes.Internal, codes.Unimplemented, codes.Unavailable:
				hub.CaptureException(err)
			}
		}()
		return handler(ctx, req)
	}
}
