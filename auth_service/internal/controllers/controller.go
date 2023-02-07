package controllers

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/ulo"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type AccountController interface {
	CreateAccount(ctx context.Context, request *proto.Credential) (
		accountID *proto.AccountID, logstash ulo.LogStash, err error,
	)
	GetToken(ctx context.Context, request *proto.Credential) (
		token *proto.Token, logstash ulo.LogStash, err error,
	)
	GetTokenPayload(ctx context.Context, request *proto.Token) (
		tokenPayload *proto.TokenPayload, logstash ulo.LogStash, err error,
	)
	Ping(ctx context.Context, request *proto.Void) (
		pong *proto.Pong, logstash ulo.LogStash, err error,
	)
}
