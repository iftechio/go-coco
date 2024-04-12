package xgrpc

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	xLogger "github.com/iftechio/go-coco/utils/logger"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/structpb"
)

type CustomError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string  `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Toast   *string `protobuf:"bytes,2,opt,name=toast,proto3,oneof" json:"toast,omitempty"`
	// E Code
	Code  *string         `protobuf:"bytes,3,opt,name=code,proto3,oneof" json:"code,omitempty"`
	Error *string         `protobuf:"bytes,4,opt,name=error,proto3,oneof" json:"error,omitempty"`
	Data  *structpb.Value `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
}

func GatewayMiddleware(
	register func(
		ctx context.Context,
		mux *runtime.ServeMux,
		endpoint string,
		opts []grpc.DialOption,
	) (err error),
	endpoint string,
	pattern string,
) echo.MiddlewareFunc {
	gateway := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
			switch s {
			case "Connection":
				return "", false
			default:
				return s, true
			}
		}),
		runtime.WithErrorHandler(gatewayErrorHandler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)
	if err := register(
		context.Background(), gateway, endpoint, []grpc.DialOption{grpc.WithInsecure()},
	); err != nil {
		panic(err)
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Path(), pattern) {
				gateway.ServeHTTP(c.Response(), c.Request())
				return nil
			}
			return next(c)
		}
	}
}

// gatewayErrorHandler 自定义 error handler
// 覆盖 runtime.DefaultErrorHandle 默认的错误处理，实现自定义的 error payload
func gatewayErrorHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	marshaler runtime.Marshaler,
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	var statusCode int
	var customStatus *runtime.HTTPStatusError
	if errors.As(err, &customStatus) {
		err = customStatus.Err
		statusCode = customStatus.HTTPStatus
	}

	s := status.Convert(err)

	pb := &CustomError{
		Message: s.Message(),
	}

	for _, d := range s.Details() {
		m, ok := d.(*CustomError)
		if !ok {
			continue
		}
		pb.Message = m.Message
		pb.Error = m.Toast
		pb.Code = m.Code
		pb.Data = m.Data
	}

	w.Header().Del("Trailer")
	w.Header().Del("Transfer-Encoding")
	w.Header().Set("Content-Type", marshaler.ContentType(pb))

	buf, merr := marshaler.Marshal(pb)
	if merr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		xLogger.F(ctx).Errorf("Failed to marshal error message %q: %v", s, merr)
		return
	}

	if md, ok := runtime.ServerMetadataFromContext(ctx); ok {
		for k, vs := range md.HeaderMD {
			h := fmt.Sprintf("%s%s", runtime.MetadataHeaderPrefix, k)
			for _, v := range vs {
				w.Header().Add(h, v)
			}
		}
	}

	if statusCode == 0 {
		statusCode = runtime.HTTPStatusFromCode(s.Code())
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(buf); err != nil {
		xLogger.F(ctx).Errorf("Failed to write response: %v", err)
	}
}
