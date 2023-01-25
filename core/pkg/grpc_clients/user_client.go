package grpcclients

import (
	"context"
	"time"

	proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ProtobufUserClient struct {
	proto.UserClient
}

func NewProtobufUserClient(ctx context.Context, serviceAddress string, connectionTimeout time.Duration,
	backoffConfig backoff.Config, keepaliveParams keepalive.ClientParameters,
	unaryInterceptors []grpc.UnaryClientInterceptor,
) (*ProtobufUserClient, error) {
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepaliveParams),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			middleware.ChainUnaryClient(
				unaryInterceptors...,
			),
		),
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
	return &ProtobufUserClient{
		UserClient: proto.NewUserClient(conn),
	}, nil
}
