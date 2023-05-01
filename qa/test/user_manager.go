package test

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type UserManager struct {
	client UserServiceHTTPClient
}

func (u *UserManager) NewUnauthorizedUser(ctx context.Context) (*UserID, *Credential, error) {
	userCredential := &Credential{
		Email:    uuid.NewString(),
		Password: uuid.NewString(),
	}

	userID, err := u.client.CreateUser(&CreateUserRequest{
		Credential: userCredential,
		UserData: &UserData{
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

func (u *UserManager) NewAuthorizedUser(ctx context.Context) (*UserID, *Token, error) {
	userID, credential, err := u.NewUnauthorizedUser(ctx)
	if err != nil {
		return nil, nil, err
	}

	token, err := u.client.GenerateToken(credential)
	if err != nil {
		return nil, nil, err
	}

	return userID, token, nil
}

func (u *UserManager) NewContextWithToken(ctx context.Context, token *Token) context.Context {
	metadataPairs := metadata.Pairs(grpcAuthorizationHeader, token.AccessToken)
	return metadata.NewOutgoingContext(ctx, metadataPairs)
}
