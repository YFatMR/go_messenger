package repositories

import "errors"

var (
	UserNotFoundErr      = errors.New("document not found")
	WrongUserIdFormatErr = errors.New("wrong format of user id")
)
