package token

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/entities"
)

type Entity struct {
	accessToken string
}

func New(accessToken string) *Entity {
	return &Entity{
		accessToken: accessToken,
	}
}

func FromProtobuf(token *proto.Token) (*Entity, error) {
	if token == nil || token.GetAccessToken() == "" {
		return nil, entities.ErrWrongRequestFormat
	}
	return New(token.GetAccessToken()), nil
}

func (e *Entity) GetAccessToken() string {
	return e.accessToken
}
