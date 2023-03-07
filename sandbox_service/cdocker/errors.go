package cdocker

import "errors"

var (
	ErrRunCode   = errors.New("can not run code")
	ErrHugeInput = errors.New("too huge file input")
)
