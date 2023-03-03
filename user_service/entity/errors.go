package entity

import "errors"

var (
	ErrWrongRequestFormat = errors.New("wrong request format")
	ErrUndefinedRole      = errors.New("undefined user role")
)
