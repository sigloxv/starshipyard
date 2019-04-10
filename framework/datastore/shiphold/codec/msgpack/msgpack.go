package msgpack

import (
	msgpack "github.com/vmihailenco/msgpack"
)

const name = "msgpack"

var Codec = new(msgpackCodec)

type msgpackCodec int

func (m msgpackCodec) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (m msgpackCodec) Unmarshal(b []byte, v interface{}) error {
	return msgpack.Unmarshal(b, v)
}

func (m msgpackCodec) Name() string {
	return name
}
