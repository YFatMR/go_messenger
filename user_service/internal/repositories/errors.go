package repositories

import "errors"

var (
	ErrUserNotFound      = errors.New("document not found")
	ErrWrongUserIDFormat = errors.New("wrong format of user id")
)
