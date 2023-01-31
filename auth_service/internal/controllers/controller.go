package controllers

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type AccountController interface {
	CreateAccount(ctx context.Context, request *proto.Credential) (
		accountID *proto.AccountID, lerr logerr.Error,
	)
	GetToken(ctx context.Context, request *proto.Credential) (
		token *proto.Token, lerr logerr.Error,
	)
	GetTokenPayload(ctx context.Context, request *proto.Token) (
		tokenPayload *proto.TokenPayload, lerr logerr.Error,
	)
	Ping(ctx context.Context, request *proto.Void) (pong *proto.Pong, lerr logerr.Error)
}
