package accountservice

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/auth/jwtmanager"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/credential"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/token"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/tokenpayload"
	"github.com/YFatMR/go_messenger/auth_service/internal/repositories"
	"github.com/YFatMR/go_messenger/auth_service/internal/services"
	"github.com/YFatMR/go_messenger/core/pkg/ulo"
)

type AccountService struct {
	accountRepository repositories.AccountRepository
	authManager       jwtmanager.Manager
}

func New(repository repositories.AccountRepository, authManager jwtmanager.Manager) *AccountService {
	return &AccountService{
		accountRepository: repository,
		authManager:       authManager,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, credential *credential.Entity) (
	*accountid.Entity, ulo.LogStash, error,
) {
	return s.accountRepository.CreateAccount(ctx, credential, entities.UserRole)
}

func (s *AccountService) GetToken(ctx context.Context, credential *credential.Entity) (
	*token.Entity, ulo.LogStash, error,
) {
	tokenPayload, hashedPassword, _, err := s.accountRepository.GetTokenPayloadWithHashedPasswordByLogin(
		ctx, credential.GetLogin(),
	)
	if err != nil {
		return nil, nil, err
	}
	if err := credential.VerifyPassword(hashedPassword); err != nil {
		return nil, ulo.FromErrorWithMsg("Can't verify password", err), services.ErrWrongCredential
	}

	token, _, err := s.authManager.GenerateToken(ctx, tokenPayload)
	if err != nil {
		return nil, nil, err
	}

	return token, nil, nil
}

func (s *AccountService) GetTokenPayload(ctx context.Context, token *token.Entity) (
	*tokenpayload.Entity, ulo.LogStash, error,
) {
	claims, _, err := s.authManager.VerifyToken(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	return claims.GetTokenPayload(), nil, nil
}
