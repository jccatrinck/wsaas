package json

import (
	"bytes"
	"unicode"
)

const (
	arrayOpenTag  = byte('[')
	arrayCloseTag = byte(']')
)

func EnsureArray(data []byte) []byte {
	data = bytes.TrimLeftFunc(data, unicode.IsSpace)

	if data[0] != arrayOpenTag {
		data = append([]byte{arrayOpenTag}, append(data, arrayCloseTag)...)
	}

	return data
}
