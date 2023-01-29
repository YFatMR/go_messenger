package cerrors

type Error interface {
	// error for logging, metrics, e.t.c not for users
	GetInternalErrorMessage() string
	// error for logging, metrics, e.t.c not for users
	GetInternalError() error
	// error for user
	GetAPIError() error
}

type customError struct {
	internalErrMessage string
	internalErr        error
	apiErr             error
}

func New(internalErrMessage string, internalErr error, apiErr error) Error {
	if internalErr == nil || apiErr == nil {
		panic("Both errors should be provide")
	}
	return &customError{
		internalErrMessage: internalErrMessage,
		internalErr:        internalErr,
		apiErr:             apiErr,
	}
}

func (e *customError) GetInternalError() error {
	if e == nil {
		return nil
	}
	return e.internalErr
}

func (e *customError) GetInternalErrorMessage() string {
	if e == nil {
		return ""
	}
	return e.internalErrMessage
}

func (e *customError) GetAPIError() error {
	if e == nil {
		return nil
	}
	return e.apiErr
}
