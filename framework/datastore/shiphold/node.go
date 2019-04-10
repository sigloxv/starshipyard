package shiphold

import (
	codec "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/codec"

	bolt "go.etcd.io/bbolt"
)

type Node interface {
	Tx
	TypeStore
	KeyValueStore
	BucketScanner
	From(addend ...string) Node
	Bucket() []string
	GetBucket(tx *bolt.Tx, children ...string) *bolt.Bucket
	CreateBucketIfNotExists(tx *bolt.Tx, bucket string) (*bolt.Bucket, error)
	WithTransaction(tx *bolt.Tx) Node
	Begin(writable bool) (Node, error)
	Codec() codec.MarshalUnmarshaler
	WithCodec(codec codec.MarshalUnmarshaler) Node
	WithBatch(enabled bool) Node
}

type node struct {
	s          *DB
	rootBucket []string
	tx         *bolt.Tx
	codec      codec.MarshalUnmarshaler
	batchMode  bool
}

func (self node) From(addend ...string) Node {
	self.rootBucket = append(self.rootBucket, addend...)
	return &self
}

func (self node) WithTransaction(tx *bolt.Tx) Node {
	self.tx = tx
	return &self
}

func (self node) WithCodec(codec codec.MarshalUnmarshaler) Node {
	self.codec = codec
	return &self
}

func (self node) WithBatch(enabled bool) Node {
	self.batchMode = enabled
	return &self
}

func (self *node) Bucket() []string {
	return self.rootBucket
}

func (self *node) Codec() codec.MarshalUnmarshaler {
	return self.codec
}

func (self *node) readWriteTx(fn func(tx *bolt.Tx) error) error {
	if self.tx != nil {
		return fn(self.tx)
	}
	if self.batchMode {
		return self.s.Bolt.Batch(func(tx *bolt.Tx) error {
			return fn(tx)
		})
	}
	return self.s.Bolt.Update(func(tx *bolt.Tx) error {
		return fn(tx)
	})
}

func (self *node) readTx(fn func(tx *bolt.Tx) error) error {
	if self.tx != nil {
		return fn(self.tx)
	}
	return self.s.Bolt.View(func(tx *bolt.Tx) error {
		return fn(tx)
	})
}
