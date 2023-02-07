package repositories

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/credential"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/tokenpayload"
	"github.com/YFatMR/go_messenger/core/pkg/ulo"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, credential *credential.Entity, role entities.Role) (
		accountID *accountid.Entity, logtash ulo.LogStash, err error,
	)
	GetTokenPayloadWithHashedPasswordByLogin(ctx context.Context, login string) (
		tokenPayload *tokenpayload.Entity, hashedPassword string, logtash ulo.LogStash, err error,
	)
}
