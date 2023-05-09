package dialog

import "errors"

// controller errors
var ErrWrongRequestFormat = errors.New("wrong request format")

// repository errors
var (
	ErrCreateDialog   = errors.New("can not create dialog")
	ErrCreateMessage  = errors.New("can not create message")
	ErrParseRequest   = errors.New("can not parse request")
	ErrGetReadMessage = errors.New("can not get last read message")
)

// kafka errors
var (
	ErrMessageCreation = errors.New("unable to create message")
	ErrMessageWriting  = errors.New("can not write message to brocker")
)

// model error
var (
	ErrFobidden = errors.New("fobidden")
)
