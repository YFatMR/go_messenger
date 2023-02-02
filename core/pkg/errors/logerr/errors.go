package logerr

import (
	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr/logerrcore"
)

type Error interface {
	// public error
	GetAPIError() error
	// check error
	HasError() bool
	// Call for library integration: message for logging
	GetLogMessage() string
	// Call for library integration: get log level
	GetLogLevel() Loglevel
	// Call for library integration: logging fields
	GetFields() []logerrcore.Field
	// Call for library integration:
	StopLogMessage()
	// Call for library integration:
	IsLogMessage() bool
}

type Loglevel uint8

const (
	UndefinedLevel Loglevel = iota
	DebugLevel
	InfoLevel
	ErrorLevel
)

type customError struct {
	loggingMessage bool
	apiErr         error
	logMessage     string
	logLevel       Loglevel
	fields         []logerrcore.Field
}

func NewError(apiErr error, logMessage string, fields ...logerrcore.Field) *customError {
	return &customError{
		loggingMessage: true,
		apiErr:         apiErr,
		logMessage:     logMessage,
		logLevel:       ErrorLevel,
		fields:         fields,
	}
}

func NewInfo(logMessage string, fields ...logerrcore.Field) *customError {
	return &customError{
		loggingMessage: true,
		logMessage:     logMessage,
		logLevel:       InfoLevel,
		fields:         fields,
	}
}

func NewDebug(logMessage string, fields ...logerrcore.Field) *customError {
	return &customError{
		loggingMessage: true,
		logMessage:     logMessage,
		logLevel:       DebugLevel,
		fields:         fields,
	}
}

func (e *customError) GetAPIError() error {
	if e == nil {
		return nil
	}
	return e.apiErr
}

func (e *customError) GetLogMessage() string {
	if e == nil {
		return ""
	}
	return e.logMessage
}

func (e *customError) GetFields() []logerrcore.Field {
	if e == nil {
		return []logerrcore.Field{}
	}
	return e.fields
}

func (e *customError) GetLogLevel() Loglevel {
	if e == nil {
		return UndefinedLevel
	}
	return e.logLevel
}

func (e *customError) HasError() bool {
	if e == nil {
		return false
	}
	return e.apiErr != nil
}

func (e *customError) StopLogMessage() {
	if e == nil {
		return
	}
	e.loggingMessage = false
}

func (e *customError) IsLogMessage() bool {
	if e == nil {
		return false
	}
	return e.loggingMessage
}
