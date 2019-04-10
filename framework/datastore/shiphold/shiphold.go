package shiphold

import (
	"bytes"
	"encoding/binary"
	"time"

	codec "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/codec"
	json "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/codec/json"

	bolt "go.etcd.io/bbolt"
)

const (
	dbinfo         = "__shiphold_db"
	metadataBucket = "__shiphold_metadata"
)

var defaultCodec = json.Codec

// TODO: Lets simplify initializaiton using smart/secure defaults, then enable
// editing default config via chaining extra options into the initializaiton
func Open(path string, options ...func(*Options) error) (*DB, error) {
	var err error
	var opts Options
	for _, option := range options {
		if err = option(&opts); err != nil {
			return nil, err
		}
	}

	s := DB{
		Bolt: opts.bolt,
	}
	n := node{
		s:          &s,
		codec:      opts.codec,
		batchMode:  opts.batchMode,
		rootBucket: opts.rootBucket,
	}

	if n.codec == nil {
		n.codec = defaultCodec
	}

	if opts.boltMode == 0 {
		opts.boltMode = 0600
	}

	if opts.boltOptions == nil {
		opts.boltOptions = &bolt.Options{Timeout: 1 * time.Second}
	}

	s.Node = &n

	// skip if UseDB option is used
	if s.Bolt == nil {
		s.Bolt, err = bolt.Open(path, opts.boltMode, opts.boltOptions)
		if err != nil {
			return nil, err
		}
	}
	return &s, nil
}

type DB struct {
	Node
	Bolt *bolt.DB
}

func (self *DB) Close() error {
	return self.Bolt.Close()
}

func toBytes(key interface{}, codec codec.MarshalUnmarshaler) ([]byte, error) {
	if key == nil {
		return nil, nil
	}
	switch t := key.(type) {
	case []byte:
		return t, nil
	case string:
		return []byte(t), nil
	case int:
		return numberToBytes(int64(t))
	case uint:
		return numberToBytes(uint64(t))
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		return numberToBytes(t)
	default:
		return codec.Marshal(key)
	}
}

func numberToBytes(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func bytesToNumber(raw []byte) (to int64, err error) {
	r := bytes.NewReader(raw)
	err = binary.Read(r, binary.BigEndian, &to)
	if err != nil {
		return 0, err
	}
	return to, nil
}
