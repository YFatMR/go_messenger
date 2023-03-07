package entity

import "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type Token struct {
	AccessToken string
}

func TokenFromString(accessToken string) *Token {
	return &Token{AccessToken: accessToken}
}

func TokenToProtobuf(token *Token) *proto.Token {
	return &proto.Token{
		AccessToken: token.AccessToken,
	}
}
