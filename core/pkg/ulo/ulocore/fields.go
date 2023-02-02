package ulocore

type Field struct {
	Key       string
	Type      FieldType
	Integer   int64
	String    string
	Interface interface{}
}

func Int(key string, value int) Field {
	return Int64(key, int64(value))
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Type: Int64Type, Integer: value}
}

func String(key string, value string) Field {
	return Field{Key: key, Type: StringType, String: value}
}

func Message(value string) Field {
	return String("msg", value)
}

func Error(err error) Field {
	return NamedError("error", err)
}

func NamedError(key string, err error) Field {
	if err == nil {
		return Skip()
	}
	return Field{Key: key, Type: ErrorType, Interface: err}
}

func Skip() Field {
	return Field{Type: SkipType}
}
