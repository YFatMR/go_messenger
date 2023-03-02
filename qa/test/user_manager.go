package test

import (
	"context"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type UserManager struct{}

func (u *UserManager) NewUnauthorizedUser(ctx context.Context) (*proto.UserID, *proto.Credential, error) {
	userCredential := &proto.Credential{
		Login:    uuid.NewString(),
		Password: uuid.NewString(),
		Role:     "user",
	}

	userID, err := frontServicegRPCClient.CreateUser(ctx, &proto.CreateUserRequest{
		Credential: userCredential,
		UserData: &proto.UserData{
			Nickname: uuid.NewString(),
			Name:     uuid.NewString(),
			Surname:  uuid.NewString(),
		},
	})
	if err != nil {
		return nil, nil, err
	}
	return userID, userCredential, nil
}

func (u *UserManager) NewAuthorizedUser(ctx context.Context) (*proto.UserID, *proto.Token, error) {
	userID, credential, err := u.NewUnauthorizedUser(ctx)
	if err != nil {
		return nil, nil, err
	}

	token, err := frontServicegRPCClient.GenerateToken(ctx, credential)
	if err != nil {
		return nil, nil, err
	}

	return userID, token, nil
}

func (u *UserManager) NewContextWithToken(ctx context.Context, token *proto.Token) context.Context {
	metadataPairs := metadata.Pairs(grpcAuthorizationHeader, token.GetAccessToken())
	return metadata.NewOutgoingContext(ctx, metadataPairs)
}
