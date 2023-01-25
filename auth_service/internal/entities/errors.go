package entities

import "errors"

var (
	ErrWrongRequestFormat = errors.New("wrong request format")
	ErrUndefinedUserRole  = errors.New("undefined user role")
)
