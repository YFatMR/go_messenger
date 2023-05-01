package dialog

import "errors"

// controller errors
var ErrWrongRequestFormat = errors.New("wrong request format")

// repository errors
var ErrCreateDialog = errors.New("can not create dialog")
var ErrCreateMessage = errors.New("can not create message")
var ErrParseRequest = errors.New("can not parse request")
