package controllers

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type UserController interface {
	Create(ctx context.Context, request *proto.CreateUserDataRequest) (userID *proto.UserID, lerr logerr.Error)
	GetByID(ctx context.Context, request *proto.UserID) (userData *proto.UserData, lerr logerr.Error)
	DeleteByID(ctx context.Context, request *proto.UserID) (void *proto.Void, lerr logerr.Error)
	Ping(ctx context.Context, request *proto.Void) (pong *proto.Pong, lerr logerr.Error)
}
