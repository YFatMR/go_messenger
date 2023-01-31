package jwtmanager

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/auth_service/internal/auth"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/token"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/tokenpayload"

	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
	"github.com/golang-jwt/jwt/v4"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	AccountID string
	UserRole  entities.Role
}

func (c *TokenClaims) GetTokenPayload() *tokenpayload.Entity {
	return tokenpayload.New(c.AccountID, c.UserRole)
}

type Manager interface {
	GenerateToken(ctx context.Context, payload *tokenpayload.Entity) (token1 *token.Entity, lerr logerr.Error)
	VerifyToken(ctx context.Context, accessToken *token.Entity) (tokenClaims *TokenClaims, lerr logerr.Error)
}

type JWTManager struct {
	secretKey               string
	tokenExpirationDuration time.Duration
	signingMethod           jwt.SigningMethod
}

func New(secretKey string, tokenExpirationDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:               secretKey,
		tokenExpirationDuration: tokenExpirationDuration,
		signingMethod:           jwt.SigningMethodHS256,
	}
}

func (m *JWTManager) GenerateToken(ctx context.Context, payload *tokenpayload.Entity) (*token.Entity, logerr.Error) {
	if payload == nil {
		return nil, logerr.NewError(auth.ErrTokenGenerationFailed, "null payload got")
	}
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenExpirationDuration)),
		},
		AccountID: payload.GetAccountID(),
		UserRole:  payload.GetUserRole(),
	}

	resultToken := jwt.NewWithClaims(m.signingMethod, claims)
	accessToken, err := resultToken.SignedString([]byte(m.secretKey))
	if err != nil {
		return nil, logerr.NewError(auth.ErrTokenGenerationFailed, "can't generate signed string", logerr.Err(err))
	}
	return token.New(accessToken), nil
}

// Check token expiration withount direct checks.
func (m *JWTManager) VerifyToken(ctx context.Context, accessToken *token.Entity) (*TokenClaims, logerr.Error) {
	if accessToken == nil {
		return nil, logerr.NewError(auth.ErrInvalidAccessToken, "null accessToken got")
	}
	token, err := jwt.ParseWithClaims(
		accessToken.GetAccessToken(),
		&TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, auth.ErrInvalidAccessToken
			}
			return []byte(m.secretKey), nil
		},
	)
	if err != nil {
		return nil, logerr.NewError(auth.ErrInvalidAccessToken, "invalid access token got", logerr.Err(err))
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, logerr.NewError(auth.ErrInvalidAccessToken, "invalid token claims")
	}
	return claims, nil
}
