package grpcclients

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"google.golang.org/grpc"
)

func NewGRPCDialogClient(ctx context.Context, serviceAddress string, connectionTimeout time.Duration,
	opts []grpc.DialOption,
) (
	proto.DialogServiceClient, error,
) {
	ctx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serviceAddress, opts...)
	if err != nil {
		return nil, err
	}
	return proto.NewDialogServiceClient(conn), nil
}
