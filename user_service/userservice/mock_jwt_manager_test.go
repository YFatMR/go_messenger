package userservice_test

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/jwtmanager"
)

type GenerateAccessTokenResponseData struct {
	AccessToken string
	Error       error
}

type VerifyTokenResponseData struct {
	TokenClaims *jwtmanager.TokenClaims
	Error       error
}

type MockJWTManager struct {
	GenerateAccessTokenResponse GenerateAccessTokenResponseData
	VerifyTokenResponse         VerifyTokenResponseData
}

func (m *MockJWTManager) GenerateAccessToken(ctx context.Context, payload jwtmanager.TokenPayload) (
	/*accessToken*/ string, error,
) {
	return m.GenerateAccessTokenResponse.AccessToken, m.GenerateAccessTokenResponse.Error
}

func (m *MockJWTManager) VerifyToken(ctx context.Context, accessToken string) (
	*jwtmanager.TokenClaims, error,
) {
	return m.VerifyTokenResponse.TokenClaims, m.VerifyTokenResponse.Error
}
