package credential

import (
	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
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

type Entity struct {
	login          string
	password       string
	hashedPassword string
	verifier       verifier
}

func New(login string, password string, hashedPassword string, verifier verifier) *Entity {
	return &Entity{
		login:          login,
		password:       password,
		hashedPassword: hashedPassword,
		verifier:       verifier,
	}
}

func FromProtobuf(credential *proto.Credential, validator passwordValidator) (
	*Entity, error,
) {
	if credential == nil || credential.GetLogin() == "" || credential.GetPassword() == "" {
		return nil, entities.ErrWrongRequestFormat
	}
	hashedPassword, err := validator.GetHasher()(credential.GetPassword())
	if err != nil {
		return nil, err
	}
	return New(credential.GetLogin(), credential.GetPassword(), hashedPassword, validator.GetVerifier()), nil
}

func (e *Entity) GetLogin() string {
	return e.login
}

func (e *Entity) GetHashedPassword() string {
	return e.hashedPassword
}

func (e *Entity) VerifyPassword(hashedPassword string) error {
	return e.verifier(hashedPassword, e.password)
}
