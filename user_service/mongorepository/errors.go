package mongorepository

import "errors"

var (
	ErrUserCreation     = errors.New("can't create user")
	ErrInternalDatabase = errors.New("internal database error")
	ErrUserNotFound     = errors.New("document not found")
	ErrGetUser          = errors.New("can't get user")
	ErrUserDeletion     = errors.New("can't delete user")
	ErrGetToken         = errors.New("can't get token")
)
