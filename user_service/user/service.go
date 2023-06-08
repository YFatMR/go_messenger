package user

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/jwtmanager"
	"github.com/YFatMR/go_messenger/user_service/apientity"
	"github.com/YFatMR/go_messenger/user_service/entity"
	"go.uber.org/zap"
)

type userService struct {
	userRepository  apientity.UserRepository
	logger          *czap.Logger
	passwordManager apientity.PasswordManager
	jwtManager      jwtmanager.Manager
}

func NewService(repository apientity.UserRepository, passwordManager apientity.PasswordManager,
	jwtManager jwtmanager.Manager, logger *czap.Logger,
) apientity.UserService {
	return &userService{
		userRepository:  repository,
		logger:          logger,
		passwordManager: passwordManager,
		jwtManager:      jwtManager,
	}
}

func (s *userService) Create(ctx context.Context, user *entity.User, usafeCredential *entity.UnsafeCredential) (
	*entity.UserID, error,
) {
	hashedPassword, err := s.passwordManager.HashPassword(usafeCredential.Password)
	if err != nil {
		s.logger.ErrorContext(ctx, "unable to hash password")
		return nil, ErrCreateUser
	}
	credential := entity.CredentialFromUnsafeCredential(usafeCredential, hashedPassword)
	return s.userRepository.Create(ctx, user, credential)
}

func (s *userService) GetByID(ctx context.Context, userID *entity.UserID) (
	*entity.User, error,
) {
	return s.userRepository.GetByID(ctx, userID)
}

func (s *userService) DeleteByID(ctx context.Context, userID *entity.UserID) error {
	return s.userRepository.DeleteByID(ctx, userID)
}

func (s *userService) GenerateToken(ctx context.Context, unsafeCredential *entity.UnsafeCredential) (
	*entity.Token, error,
) {
	account, err := s.userRepository.GetAccountByEmail(ctx, unsafeCredential.Email)
	if err != nil {
		return nil, err
	}
	if err = s.passwordManager.VerifyPassword(account.HashedPassword, unsafeCredential.Password); err != nil {
		s.logger.ErrorContext(ctx, "wrong password", zap.String("email", unsafeCredential.Email))
		return nil, ErrWrongCredential
	}

	payload := TokenPayloadFromAccount(account)
	accessToken, err := s.jwtManager.GenerateAccessToken(ctx, payload)
	if err != nil {
		s.logger.ErrorContext(ctx, "can't generate access token", zap.String("email", unsafeCredential.Email))
		return nil, err
	}
	return entity.TokenFromString(accessToken), nil
}

func (s *userService) UpdateUserData(ctx context.Context, userID *entity.UserID,
	request *entity.UpdateUserRequest,
) error {
	return s.userRepository.UpdateUserData(ctx, userID, request)
}

func (s *userService) GetUsersByPrefix(ctx context.Context, selfID *entity.UserID, prefix string, limit uint64) (
	[]*entity.UserWithID, error,
) {
	return s.userRepository.GetUsersByPrefix(ctx, selfID, prefix, limit)
}
