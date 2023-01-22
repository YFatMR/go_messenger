package grpcclients

import (
	"context"
	"time"

	proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ProtobufAuthClient struct {
	proto.AuthClient
}

func NewProtobufAuthClient(ctx context.Context, serviceAddress string, connectionTimeout time.Duration,
	backoffConfig backoff.Config, keepaliveParams keepalive.ClientParameters,
) (*ProtobufAuthClient, error) {
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepaliveParams),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithConnectParams(
			grpc.ConnectParams{
				Backoff: backoffConfig,
			},
		),
	}

	ctx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serviceAddress, opts...)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return &ProtobufAuthClient{
		AuthClient: proto.NewAuthClient(conn),
	}, nil
}
