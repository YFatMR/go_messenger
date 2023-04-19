package grpcapi

import "errors"

var (
	ErrNoMetadata                       = errors.New("expected metadata from call")
	ErrUnexpectedMetadataAccountIDCount = errors.New("please, provide only one accountID")
)
