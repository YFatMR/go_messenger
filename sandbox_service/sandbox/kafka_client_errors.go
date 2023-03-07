package sandbox

import "errors"

var (
	ErrMessageCreation = errors.New("unable to create message")
	ErrMessageWriting  = errors.New("can not write message to brocker")
)
