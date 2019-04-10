package codec

type Codec int

const (
	Gob Codec = iota
	JSON
	MsgPack
)

type MarshalUnmarshaler interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(b []byte, v interface{}) error
	Name() string
}
