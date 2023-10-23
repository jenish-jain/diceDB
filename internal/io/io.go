package io

type Decoder interface {
	Decode(data []byte) (interface{}, error)
}

type Encoder interface {
	Encode(value interface{}, isSimple bool) []byte
}
