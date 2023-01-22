package entities

import proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type Token struct {
	accessToken string
}

func newToken(accessToken string) *Token {
	return &Token{
		accessToken: accessToken,
	}
}

func NewTokenFromProtobuf(token *proto.Token) (*Token, error) {
	if token == nil || token.GetAccessToken() == "" {
		return nil, ErrWrongRequestFormat
	}
	return newToken(token.GetAccessToken()), nil
}

func NewTokenFromRawTokenClaims(accessToken string) *Token {
	return newToken(accessToken)
}

func (t *Token) GetAccessToken() string {
	return t.accessToken
}
