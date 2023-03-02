package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

// UnsafeProxyController provide API for handlers without authorization.
type UnsafeProxyController interface {
	CreateUser(ctx context.Context, request *proto.CreateUserRequest) (
		userID *proto.UserID, err error,
	)
	GenerateToken(ctx context.Context, request *proto.Credential) (
		token *proto.Token, err error,
	)
	Ping(ctx context.Context, request *proto.Void) (
		pong *proto.Pong, err error,
	)
}
