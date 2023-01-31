package servers

import (
	"context"

	"github.com/YFatMR/go_messenger/auth_service/internal/controllers"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type GRPCAuthServer struct {
	proto.UnimplementedAuthServer
	accountController controllers.AccountController
}

func NewGRPCAuthServer(accountController controllers.AccountController) GRPCAuthServer {
	return GRPCAuthServer{
		accountController: accountController,
	}
}

func (s *GRPCAuthServer) CreateAccount(ctx context.Context, request *proto.Credential) (
	*proto.AccountID, error,
) {
	accountID, lerr := s.accountController.CreateAccount(ctx, request)
	return accountID, lerr.GetAPIError()
}

func (s *GRPCAuthServer) GetToken(ctx context.Context, request *proto.Credential) (*proto.Token, error) {
	token, lerr := s.accountController.GetToken(ctx, request)
	return token, lerr.GetAPIError()
}

func (s *GRPCAuthServer) GetTokenPayload(ctx context.Context, request *proto.Token) (*proto.TokenPayload, error) {
	tokenPayload, lerr := s.accountController.GetTokenPayload(ctx, request)
	return tokenPayload, lerr.GetAPIError()
}

func (s *GRPCAuthServer) Ping(ctx context.Context, request *proto.Void) (*proto.Pong, error) {
	pong, lerr := s.accountController.Ping(ctx, request)
	return pong, lerr.GetAPIError()
}
