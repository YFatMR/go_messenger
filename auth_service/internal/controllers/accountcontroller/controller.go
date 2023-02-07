package accountcontroller

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/controllers"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/credential"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/token"
	"github.com/YFatMR/go_messenger/auth_service/internal/services"
	"github.com/YFatMR/go_messenger/core/pkg/ulo"
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
	*proto.AccountID, ulo.LogStash, error,
) {
	credential, err := credential.FromProtobuf(request, c.passwordValidator)
	if err != nil {
		logstash := ulo.FromErrorWithMsg("Can't parse credential", err)
		return nil, logstash, controllers.ErrWrongRequestFormat
	}

	accountID, _, err := c.accountService.CreateAccount(ctx, credential)
	if err != nil {
		return nil, nil, err
	}

	return &proto.AccountID{
		ID: accountID.GetID(),
	}, nil, nil
}

func (c *AccountController) GetToken(ctx context.Context, request *proto.Credential) (
	*proto.Token, ulo.LogStash, error,
) {
	credential, err := credential.FromProtobuf(request, c.passwordValidator)
	if err != nil {
		logstash := ulo.FromErrorWithMsg("Can't parse credential", err)
		return nil, logstash, controllers.ErrWrongRequestFormat
	}

	token, _, err := c.accountService.GetToken(ctx, credential)
	if err != nil {
		return nil, nil, err
	}

	return &proto.Token{
		AccessToken: token.GetAccessToken(),
	}, nil, nil
}

func (c *AccountController) GetTokenPayload(ctx context.Context, request *proto.Token) (
	*proto.TokenPayload, ulo.LogStash, error,
) {
	token, err := token.FromProtobuf(request)
	if err != nil {
		logstash := ulo.FromErrorWithMsg("Can't parse token", err)
		return nil, logstash, controllers.ErrWrongRequestFormat
	}

	payload, _, err := c.accountService.GetTokenPayload(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	role := payload.GetUserRole()
	return &proto.TokenPayload{
		AccountID: payload.GetAccountID(),
		UserRole:  string(role),
	}, nil, nil
}

func (c *AccountController) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, ulo.LogStash, error,
) {
	return &proto.Pong{
		Message: "pong",
	}, nil, nil
}
