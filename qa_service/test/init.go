package test

import (
	"context"
	"time"

	grpcclients "github.com/YFatMR/go_messenger/core/pkg/grpc_clients"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	frontServicegRPCClient  proto.FrontClient
	grpcAuthorizationHeader string
)

func newProtobufAuthClient(ctx context.Context, serviceAddress string, responseTimeout time.Duration) (
	proto.AuthClient, error,
) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	}
	return grpcclients.NewProtobufAuthClient(ctx, serviceAddress, responseTimeout, opts)
}

func newProtobufUserClient(ctx context.Context, serviceAddress string, responseTimeout time.Duration) (
	proto.UserClient, error,
) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	}
	return grpcclients.NewProtobufUserClient(ctx, serviceAddress, responseTimeout, opts)
}

func newProtobufFrontClient(ctx context.Context, serviceAddress string, responseTimeout time.Duration) (
	proto.FrontClient, error,
) {
	ctx, cancel := context.WithTimeout(ctx, responseTimeout)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	}

	conn, err := grpc.DialContext(ctx, serviceAddress, opts...)
	if err != nil {
		return nil, err
	}

	return proto.NewFrontClient(conn), nil
}

func pingService(ctx context.Context, pingCallback func(context.Context) (*proto.Pong, error), serviceSetupRetryCount int, serviceSetupRetryInterval time.Duration) error {
	var err error
	for i := 0; i < serviceSetupRetryCount; i++ {
		_, err = pingCallback(ctx)
		if err == nil {
			return nil
		}
	}
	return err
}
