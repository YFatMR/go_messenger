package controllers

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/ulo"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type UserController interface {
	Create(ctx context.Context, request *proto.CreateUserDataRequest) (
		userID *proto.UserID, logstash ulo.LogStash, err error,
	)
	GetByID(ctx context.Context, request *proto.UserID) (
		userData *proto.UserData, logstash ulo.LogStash, err error,
	)
	DeleteByID(ctx context.Context, request *proto.UserID) (
		void *proto.Void, logstash ulo.LogStash, err error,
	)
	Ping(ctx context.Context, request *proto.Void) (
		pong *proto.Pong, logstash ulo.LogStash, err error,
	)
}
