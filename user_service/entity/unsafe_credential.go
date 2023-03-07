package entity

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type UnsafeCredential struct {
	Login    string
	Password string
	Role     *UserRole
}

func UnsafeCredentialFromProtobuf(credential *proto.Credential) (
	*UnsafeCredential, error,
) {
	if credential == nil || credential.Login == "" || credential.Password == "" || credential.Role == "" {
		return nil, ErrWrongRequestFormat
	}
	role, err := UserRoleFromString(credential.Role)
	if err != nil {
		return nil, err
	}
	return &UnsafeCredential{
		Login:    credential.Login,
		Password: credential.Password,
		Role:     role,
	}, nil
}
