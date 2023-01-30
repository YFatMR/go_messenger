package accountcontroller

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/controllers"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/credential"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/token"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/tokenpayload"
	"github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
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
	CreateAccount(context.Context, *credential.Entity) (*accountid.Entity, error)
	GetToken(context.Context, *credential.Entity) (*token.Entity, error)
	GetTokenPayload(context.Context, *token.Entity) (*tokenpayload.Entity, error)
}

type AccountController struct {
	accountService    accountService
	passwordValidator passwordValidator
}

func New(accountService accountService, passwordValidator passwordValidator) *AccountController {
	return &AccountController{
		accountService:    accountService,
		passwordValidator: passwordValidator,
	}
}

func (c *AccountController) CreateAccount(ctx context.Context, request *proto.Credential) (
	*proto.AccountID, cerrors.Error,
) {
	credential, err := entities.NewCredentialFromProtobuf(request, c.passwordValidator)
	if err != nil {
		return nil, cerrors.New("Can't parse credential", err, controllers.ErrWrongRequestFormat)
	}

	accountID, err := c.accountService.CreateAccount(ctx, credential)
	if err != nil {
		return nil, cerrors.New("Can't create account", err, err)
	}

	return &proto.AccountID{
		ID: accountID.GetID(),
	}, nil
}

func (c *AccountController) GetToken(ctx context.Context, request *proto.Credential) (
	*proto.Token, cerrors.Error,
) {
	credential, err := entities.NewCredentialFromProtobuf(request, c.passwordValidator)
	if err != nil {
		return nil, cerrors.New("Can't parse credential", err, controllers.ErrWrongRequestFormat)
	}

	token, err := c.accountService.GetToken(ctx, credential)
	if err != nil {
		return nil, cerrors.New("Can't get token", err, err)
	}

	return &proto.Token{
		AccessToken: token.GetAccessToken(),
	}, nil
}

func (c *AccountController) GetTokenPayload(ctx context.Context, request *proto.Token) (
	*proto.TokenPayload, cerrors.Error,
) {
	token, err := entities.NewTokenFromProtobuf(request)
	if err != nil {
		return nil, cerrors.New("Can't parse token", err, controllers.ErrWrongRequestFormat)
	}

	payload, err := c.accountService.GetTokenPayload(ctx, token)
	if err != nil {
		return nil, cerrors.New("Can't get token payload", err, err)
	}

	role := payload.GetUserRole()
	return &proto.TokenPayload{
		AccountID: payload.GetAccountID(),
		UserRole:  string(role),
	}, nil
}

func (c *AccountController) Ping(ctx context.Context, request *proto.Void) (*proto.Pong, cerrors.Error) {
	return &proto.Pong{
		Message: "pong",
	}, nil
}
