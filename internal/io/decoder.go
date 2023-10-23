package io

import (
	"errors"
)

// protocol spec : https://redis.io/docs/reference/protocol-spec/
type respImpl struct {
}

// reads the length typically the first integer of the string
// until hit by a non-digit byte and returns
// the integer and the delta = length + 2 (CRLF)
func getLength(data []byte) (int, int) {
	pos, length := 0, 0
	for pos = range data {
		b := data[pos]
		if !(b >= '0' && b <= '9') {
			return length, pos + 2
		}
		length = length*10 + int(b-'0')
	}
	return 0, 0
}

// reads a RESP encoded simple string from data and returns
// the string, the delta, and the error
// +<data>\r\n
func readSimpleString(data []byte) (string, int, error) {
	// first character +
	pos := 1

	for ; data[pos] != '\r'; pos++ {
	}
	return string(data[1:pos]), pos + 2, nil
}

// reads a RESP encoded error from data and returns
// the error string, the delta, and the error
// +<data>\r\n
func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

// reads a RESP encoded integer from data and returns
// the integer value, the delta, and the error
// :[<+|->]<value>\r\n
func readInt64(data []byte) (int64, int, error) {
	// first character :
	pos := 1
	var value int64 = 0

	for ; data[pos] != '\r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}

	return value, pos + 2, nil
}

// reads a RESP encoded string from data and returns
// the string, the delta, and the error
// $<length>\r\n<data>\r\n
func readBulkString(data []byte) (string, int, error) {
	// first character $
	pos := 1

	// reading the length and forwarding the pos by
	// the length of the integer + the first special character
	length, delta := getLength(data[pos:])
	pos += delta

	// reading `len` bytes as string
	return string(data[pos:(pos + length)]), pos + length + 2, nil
}

// reads a RESP encoded array from data and returns
// the array, the delta, and the error
// *<number-of-elements>\r\n<element-1>...<element-n>
func readArray(data []byte) (interface{}, int, error) {
	// first character *
	pos := 1

	length, finalPos := getLength(data[pos:])
	pos += finalPos

	var elems = make([]interface{}, length)
	for i := range elems {
		elem, delta, err := decodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elems[i] = elem
		pos += delta
	}
	return elems, pos, nil
}

func decodeOne(data []byte) (interface{}, int, error) {

	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}

	return nil, 0, nil
}

func (r respImpl) Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	value, _, err := decodeOne(data)
	return value, err
}

func NewRESPDecoder() Decoder {
	return &respImpl{}
}
