package controllers

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type (
	hasher   = func(string) (string, error)
	verifier = func(hashedPassword string, password string) error
)

type passwordValidator interface {
	GetHasher() hasher
	GetVerifier() verifier
}

type accountService interface {
	CreateAccount(context.Context, *entities.Credential) (*entities.AccountID, error)
	GetToken(context.Context, *entities.Credential) (*entities.Token, error)
	GetTokenPayload(context.Context, *entities.Token) (*entities.TokenPayload, error)
}

type AccountController struct {
	accountService    accountService
	passwordValidator passwordValidator
	logger            *loggers.OtelZapLoggerWithTraceID
	tracer            trace.Tracer
}

func NewAccountController(accountService accountService, passwordValidator passwordValidator,
	logger *loggers.OtelZapLoggerWithTraceID, tracer trace.Tracer,
) *AccountController {
	return &AccountController{
		accountService:    accountService,
		passwordValidator: passwordValidator,
		logger:            logger,
		tracer:            tracer,
	}
}

func (c *AccountController) CreateAccount(ctx context.Context, request *proto.Credential) (
	*proto.AccountID, error,
) {
	credential, err := entities.NewCredentialFromProtobuf(request, c.passwordValidator)
	if err != nil {
		c.logger.ErrorContextNoExport(ctx, "Can't parse credential", zap.Error(err))
		return nil, ErrWrongRequestFormat
	}
	c.logger.InfoContextNoExport(ctx, "Credential parsed successfully")

	accountID, err := c.accountService.CreateAccount(ctx, credential)
	if err != nil {
		c.logger.ErrorContextNoExport(ctx, "Can't create account", zap.Error(err))
		return nil, err
	}

	return &proto.AccountID{
		ID: accountID.GetID(),
	}, nil
}

func (c *AccountController) GetToken(ctx context.Context, request *proto.Credential) (
	*proto.Token, error,
) {
	credential, err := entities.NewCredentialFromProtobuf(request, c.passwordValidator)
	if err != nil {
		c.logger.ErrorContextNoExport(ctx, "Can't parse credential", zap.Error(err))
		return nil, ErrWrongRequestFormat
	}
	c.logger.InfoContextNoExport(ctx, "Credential parsed successfully")

	token, err := c.accountService.GetToken(ctx, credential)
	if err != nil {
		c.logger.ErrorContextNoExport(ctx, "Can't get token", zap.Error(err))
		return nil, err
	}

	return &proto.Token{
		AccessToken: token.GetAccessToken(),
	}, nil
}

func (c *AccountController) GetTokenPayload(ctx context.Context, request *proto.Token) (*proto.TokenPayload, error) {
	token, err := entities.NewTokenFromProtobuf(request)
	if err != nil {
		c.logger.ErrorContextNoExport(ctx, "Can't parse token", zap.Error(err))
		return nil, ErrWrongRequestFormat
	}
	c.logger.InfoContextNoExport(ctx, "Token parsed successfully")

	payload, err := c.accountService.GetTokenPayload(ctx, token)
	if err != nil {
		return nil, err
	}

	role := payload.GetUserRole()
	return &proto.TokenPayload{
		AccountID: payload.GetAccountID(),
		UserRole:  string(role),
	}, nil
}
