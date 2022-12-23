package user_server

import (
	"context"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterRestUserServer(ctx context.Context, mux *runtime.ServeMux, grpcFrontServerAddress string) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := proto.RegisterFrontUserHandlerFromEndpoint(ctx, mux, grpcFrontServerAddress, opts); err != nil {
		panic(err)
	}
}

func RegisterGrpcUserServer(grpcServer grpc.ServiceRegistrar, userServerAddress string, logger *loggers.OtelZapLoggerWithTraceID, tracer trace.Tracer) {
	proto.RegisterFrontUserServer(grpcServer, newFrontUserServer(userServerAddress, logger, tracer))
}
