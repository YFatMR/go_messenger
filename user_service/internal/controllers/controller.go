package controllers

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type UserController interface {
	Create(ctx context.Context, request *proto.CreateUserDataRequest) (userID *proto.UserID, cerr cerrors.Error)
	GetByID(ctx context.Context, request *proto.UserID) (userData *proto.UserData, cerr cerrors.Error)
	DeleteByID(ctx context.Context, request *proto.UserID) (void *proto.Void, cerr cerrors.Error)
	Ping(ctx context.Context, request *proto.Void) (pong *proto.Pong, cerr cerrors.Error)
}
