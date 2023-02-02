package ulocore

type FieldType uint8

const (
	UnknownType FieldType = iota
	StringType
	ErrorType
	IntType
	Int64Type
	BoolType
	SkipType
)

type LogLevel uint8

const (
	UnknownLevel LogLevel = iota
	DebugLevel
	InfoLevel
	WarningLevel
	ErrorLevel
)
