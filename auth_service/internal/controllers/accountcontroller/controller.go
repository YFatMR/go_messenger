package accountcontroller

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/controllers"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/credential"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/token"
	"github.com/YFatMR/go_messenger/auth_service/internal/services"

	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
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

type AccountController struct {
	accountService    services.AccountService
	passwordValidator passwordValidator
}

func New(accountService services.AccountService, passwordValidator passwordValidator) *AccountController {
	return &AccountController{
		accountService:    accountService,
		passwordValidator: passwordValidator,
	}
}

func (c *AccountController) CreateAccount(ctx context.Context, request *proto.Credential) (
	*proto.AccountID, logerr.Error,
) {
	credential, err := credential.FromProtobuf(request, c.passwordValidator)
	if err != nil {
		return nil, logerr.NewError(controllers.ErrWrongRequestFormat, "Can't parse credential", logerr.Err(err))
	}

	accountID, lerr := c.accountService.CreateAccount(ctx, credential)
	if lerr != nil {
		return nil, lerr
	}

	return &proto.AccountID{
		ID: accountID.GetID(),
	}, nil
}

func (c *AccountController) GetToken(ctx context.Context, request *proto.Credential) (
	*proto.Token, logerr.Error,
) {
	credential, err := credential.FromProtobuf(request, c.passwordValidator)
	if err != nil {
		return nil, logerr.NewError(controllers.ErrWrongRequestFormat, "Can't parse credential", logerr.Err(err))
	}

	token, lerr := c.accountService.GetToken(ctx, credential)
	if lerr != nil {
		return nil, lerr
	}

	return &proto.Token{
		AccessToken: token.GetAccessToken(),
	}, nil
}

func (c *AccountController) GetTokenPayload(ctx context.Context, request *proto.Token) (
	*proto.TokenPayload, logerr.Error,
) {
	token, err := token.FromProtobuf(request)
	if err != nil {
		return nil, logerr.NewError(controllers.ErrWrongRequestFormat, "Can't parse token", logerr.Err(err))
	}

	payload, lerr := c.accountService.GetTokenPayload(ctx, token)
	if lerr != nil {
		return nil, lerr
	}

	role := payload.GetUserRole()
	return &proto.TokenPayload{
		AccountID: payload.GetAccountID(),
		UserRole:  string(role),
	}, nil
}

func (c *AccountController) Ping(ctx context.Context, request *proto.Void) (*proto.Pong, logerr.Error) {
	return &proto.Pong{
		Message: "pong",
	}, nil
}
