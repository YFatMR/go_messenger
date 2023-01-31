package services

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/credential"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/token"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/tokenpayload"
	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
)

type AccountService interface {
	CreateAccount(ctx context.Context, credential *credential.Entity) (
		accountID *accountid.Entity, lerr logerr.Error,
	)
	GetToken(ctx context.Context, credential *credential.Entity) (token *token.Entity, lerr logerr.Error)
	GetTokenPayload(ctx context.Context, token *token.Entity) (tokenPayload *tokenpayload.Entity, lerr logerr.Error)
}
