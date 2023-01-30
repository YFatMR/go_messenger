package controllers

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type AccountController interface {
	CreateAccount(ctx context.Context, request *proto.Credential) (
		accountID *proto.AccountID, cerr cerrors.Error,
	)
	GetToken(ctx context.Context, request *proto.Credential) (
		token *proto.Token, cerr cerrors.Error,
	)
	GetTokenPayload(ctx context.Context, request *proto.Token) (
		tokenPayload *proto.TokenPayload, cerr cerrors.Error,
	)
	Ping(ctx context.Context, request *proto.Void) (pong *proto.Pong, cerr cerrors.Error)
}
