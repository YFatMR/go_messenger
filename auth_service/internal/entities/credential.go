package entities

import (
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

type Credential struct {
	login          string
	password       string
	hashedPassword string
	verifier       verifier
}

func NewCredential(login string, password string, hashedPassword string, verifier verifier) *Credential {
	return &Credential{
		login:          login,
		password:       password,
		hashedPassword: hashedPassword,
		verifier:       verifier,
	}
}

func NewCredentialFromProtobuf(credential *proto.Credential, validator passwordValidator) (
	*Credential, error,
) {
	if credential == nil || credential.GetLogin() == "" || credential.GetPassword() == "" {
		return nil, ErrWrongRequestFormat
	}
	hashedPassword, err := validator.GetHasher()(credential.GetPassword())
	if err != nil {
		return nil, err
	}
	return NewCredential(credential.GetLogin(), credential.GetPassword(), hashedPassword, validator.GetVerifier()), nil
}

func (c *Credential) GetLogin() string {
	return c.login
}

func (c *Credential) GetHashedPassword() string {
	return c.hashedPassword
}

func (c *Credential) VerifyPassword(hashedPassword string) error {
	return c.verifier(hashedPassword, c.password)
}
