package ezsec

type valueType int

const (
	valueTypeByteSlice valueType = iota
	valueTypeString
)

func valueFromInterface(value interface{}, resultType valueType) interface{} {
	var result interface{}

	switch value.(type) {
	case []byte:

		switch resultType {
		case valueTypeByteSlice:
			result = value.([]byte)
		case valueTypeString:
			result = string(value.([]byte))
		}
	case string:
		switch resultType {
		case valueTypeByteSlice:
			result = []byte(value.(string))
		case valueTypeString:
			result = value.(string)
		}
	}

	return result
}
