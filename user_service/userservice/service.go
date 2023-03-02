package userservice

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/jwtmanager"
	"github.com/YFatMR/go_messenger/user_service/entities/credential"
	"github.com/YFatMR/go_messenger/user_service/entities/token"
	"github.com/YFatMR/go_messenger/user_service/entities/unsafecredential"
	"github.com/YFatMR/go_messenger/user_service/entities/user"
	"github.com/YFatMR/go_messenger/user_service/entities/userid"
	"github.com/YFatMR/go_messenger/user_service/passwordmanager"
	"github.com/YFatMR/go_messenger/user_service/repositories"
	"github.com/YFatMR/go_messenger/user_service/services"
	"go.uber.org/zap"
)

type userService struct {
	userRepository  repositories.UserRepository
	logger          *czap.Logger
	passwordManager passwordmanager.Manager
	jwtManager      jwtmanager.Manager
}

func New(repository repositories.UserRepository, passwordManager passwordmanager.Manager,
	jwtManager jwtmanager.Manager, logger *czap.Logger,
) services.UserService {
	return &userService{
		userRepository:  repository,
		logger:          logger,
		passwordManager: passwordManager,
		jwtManager:      jwtManager,
	}
}

func (s *userService) Create(ctx context.Context, user *user.Entity, usafeCredential *unsafecredential.Entity) (
	*userid.Entity, error,
) {
	hashedPassword, err := s.passwordManager.HashPassword(usafeCredential.GetPassword())
	if err != nil {
		s.logger.ErrorContext(ctx, "unable to hash password")
		return nil, ErrCreateUser
	}
	credential := credential.New(usafeCredential.GetLogin(), hashedPassword, usafeCredential.GetRole())
	return s.userRepository.Create(ctx, user, credential)
}

func (s *userService) GetByID(ctx context.Context, userID *userid.Entity) (
	*user.Entity, error,
) {
	return s.userRepository.GetByID(ctx, userID)
}

func (s *userService) DeleteByID(ctx context.Context, userID *userid.Entity) error {
	return s.userRepository.DeleteByID(ctx, userID)
}

func (s *userService) GenerateToken(ctx context.Context, unsafeCredential *unsafecredential.Entity) (
	*token.Entity, error,
) {
	account, err := s.userRepository.GetAccountByLogin(ctx, unsafeCredential.GetLogin())
	if err != nil {
		return nil, err
	}
	if err = s.passwordManager.VerifyPassword(account.GetHashedPassword(), unsafeCredential.GetPassword()); err != nil {
		s.logger.ErrorContext(ctx, "wrong password", zap.String("login", unsafeCredential.GetLogin()))
		return nil, ErrWrongCredential
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(ctx, jwtmanager.TokenPayload{
		UserID:   account.GetUserID(),
		UserRole: account.GetRole().GetName(),
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "can't generate access token", zap.String("login", unsafeCredential.GetLogin()))
		return nil, err
	}
	return token.New(accessToken), nil
}
