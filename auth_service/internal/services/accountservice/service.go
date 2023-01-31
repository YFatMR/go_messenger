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

	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
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
	*accountid.Entity, logerr.Error,
) {
	return s.accountRepository.CreateAccount(ctx, credential, entities.UserRole)
}

func (s *AccountService) GetToken(ctx context.Context, credential *credential.Entity) (*token.Entity, logerr.Error) {
	tokenPayload, hashedPassword, lerr := s.accountRepository.GetTokenPayloadWithHashedPasswordByLogin(
		ctx, credential.GetLogin(),
	)
	if lerr != nil {
		return nil, lerr
	}

	if err := credential.VerifyPassword(hashedPassword); err != nil {
		return nil, logerr.NewError(services.ErrWrongCredential, "Can't verify password", logerr.Err(err))
	}

	token, lerr := s.authManager.GenerateToken(ctx, tokenPayload)
	if lerr != nil {
		return nil, lerr
	}

	return token, nil
}

func (s *AccountService) GetTokenPayload(ctx context.Context, token *token.Entity) (
	*tokenpayload.Entity, logerr.Error,
) {
	claims, err := s.authManager.VerifyToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return claims.GetTokenPayload(), nil
}
