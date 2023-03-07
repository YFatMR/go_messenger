package sandbox

import "errors"

var (
	ErrProgramExecution                 = errors.New("can't execute program")
	ErrNoMetadata                       = errors.New("expected metadata from call")
	ErrNoMetadataKey                    = errors.New("not found key from metadata")
	ErrUnexpectedMetadataAccountIDCount = errors.New("please, provide only one accountID")
	ErrWrongRequestFormat               = errors.New("wrong request format")
)
