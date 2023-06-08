package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type UserController interface {
	Create(ctx context.Context, request *proto.CreateUserRequest) (
		userID *proto.UserID, err error,
	)
	GetByID(ctx context.Context, request *proto.UserID) (
		userData *proto.UserData, err error,
	)
	DeleteByID(ctx context.Context, request *proto.UserID) (
		void *proto.Void, err error,
	)
	GenerateToken(ctx context.Context, request *proto.Credential) (
		void *proto.Token, err error,
	)
	UpdateUserData(ctx context.Context, request *proto.UpdateUserDataRequest) (
		void *proto.Void, err error,
	)
	GetUsersByPrefix(ctx context.Context, request *proto.GetUsersByPrefixRequest) (
		resp *proto.GetUsersByPrefixResponse, err error,
	)
	Ping(ctx context.Context, request *proto.Void) (
		pong *proto.Pong, err error,
	)
}
