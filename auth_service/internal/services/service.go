package services

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/auth"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/metrics/prometheus"
	"go.opentelemetry.io/otel/trace"
)

type accountRepository interface {
	CreateAccount(context.Context, *entities.Credential, entities.Role) (_ *entities.AccountID, err error)
	GetTokenPayloadWithHashedPasswordByLogin(context.Context, string) (
		_ *entities.TokenPayload, hashedPassword string, err error,
	)
}

type authManager interface {
	GenerateToken(ctx context.Context, payload *entities.TokenPayload) (*entities.Token, error)
	VerifyToken(ctx context.Context, accessToken *entities.Token) (*auth.TokenClaims, error)
}

type AccountService struct {
	accountRepository accountRepository
	authManager       authManager
	logger            *loggers.OtelZapLoggerWithTraceID
	tracer            trace.Tracer
}

func NewAccountService(repository accountRepository, authManager authManager,
	logger *loggers.OtelZapLoggerWithTraceID, tracer trace.Tracer,
) *AccountService {
	return &AccountService{
		accountRepository: repository,
		authManager:       authManager,
		logger:            logger,
		tracer:            tracer,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, credential *entities.Credential) (
	_ *entities.AccountID, err error,
) {
	const endpointTag = "CreateAccount"
	defer prometheus.CollectServiceRequestMetrics(endpointTag, err)

	// use default role for public api
	return s.accountRepository.CreateAccount(ctx, credential, entities.UserRole)
}

func (s *AccountService) GetToken(ctx context.Context, credential *entities.Credential) (_ *entities.Token, err error) {
	const endpointTag = "GetToken"
	defer prometheus.CollectServiceRequestMetrics(endpointTag, err)

	tokenPayload, hashedPassword, err := s.accountRepository.GetTokenPayloadWithHashedPasswordByLogin(
		ctx, credential.GetLogin(),
	)
	if err != nil {
		return nil, err
	}

	if err := credential.VerifyPassword(hashedPassword); err != nil {
		return nil, ErrWrongCredential
	}
	s.logger.InfoContextNoExport(ctx, "Password verified")

	token, err := s.authManager.GenerateToken(ctx, tokenPayload)
	if err != nil {
		return nil, err
	}

	s.logger.InfoContextNoExport(ctx, "Token generated successfully")
	return token, err
}

func (s *AccountService) GetTokenPayload(ctx context.Context, token *entities.Token) (
	_ *entities.TokenPayload, err error,
) {
	const endpointTag = "GetTokenPayload"
	defer prometheus.CollectServiceRequestMetrics(endpointTag, err)

	claims, err := s.authManager.VerifyToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return claims.GetTokenPayload(), nil
}
