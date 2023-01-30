package repositories

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrAccountCreation = errors.New("can't create account")
	ErrGetToken        = errors.New("can't get token")
)
