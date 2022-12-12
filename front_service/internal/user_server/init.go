package user_server

import (
	"context"
	. "core/pkg/loggers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "protocol/pkg/proto"
)

func RegisterRestUserServer(ctx context.Context, mux *runtime.ServeMux, grpcFrontServerAddress string) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := proto.RegisterFrontUserHandlerFromEndpoint(ctx, mux, grpcFrontServerAddress, opts)
	if err != nil {
		panic(err)
	}
}

func RegisterGrpcUserServer(grpcServer grpc.ServiceRegistrar, userServerAddress string, logger *OtelZapLoggerWithTraceID, tracer trace.Tracer) {
	proto.RegisterFrontUserServer(grpcServer, newFrontUserServer(userServerAddress, logger, tracer))
}
