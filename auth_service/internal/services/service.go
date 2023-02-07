package services

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/credential"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/token"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/tokenpayload"
	"github.com/YFatMR/go_messenger/core/pkg/ulo"
)

type AccountService interface {
	CreateAccount(ctx context.Context, credential *credential.Entity) (
		accountID *accountid.Entity, logStash ulo.LogStash, err error,
	)
	GetToken(ctx context.Context, credential *credential.Entity) (
		token *token.Entity, logStash ulo.LogStash, err error,
	)
	GetTokenPayload(ctx context.Context, token *token.Entity) (
		tokenPayload *tokenpayload.Entity, logStash ulo.LogStash, err error,
	)
}
