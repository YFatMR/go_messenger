package entities

import "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type Token struct {
	accessToken string
}

func NewToken(accessToken string) *Token {
	return &Token{
		accessToken: accessToken,
	}
}

func NewTokenFromProtobuf(token *proto.Token) (*Token, error) {
	if token == nil || token.GetAccessToken() == "" {
		return nil, ErrWrongRequestFormat
	}
	return NewToken(token.GetAccessToken()), nil
}

func (t *Token) GetAccessToken() string {
	return t.accessToken
}
