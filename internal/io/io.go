package io

type Decoder interface {
	Decode(data []byte) (interface{}, error)
}
