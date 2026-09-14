package libraries

import "bytes"

func CheckNullValue(value []byte) bool {
	return value == nil || bytes.Equal(value, []byte("null"))
}
