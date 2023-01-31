package logerrcore

type FieldType uint8

const (
	UnknownType FieldType = iota
	StringType
	Int64Type
	ErrorType
	BoolType
	SkipType
)

type Field struct {
	Key     string
	Type    FieldType
	Integer int64
	String  string
	Error   error
}
