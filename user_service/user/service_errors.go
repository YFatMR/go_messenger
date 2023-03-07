package user

import "errors"

var (
	ErrCreateUser      = errors.New("can't create user")
	ErrWrongCredential = errors.New("wrong credential")
)
