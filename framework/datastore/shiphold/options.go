package shiphold

import (
	"os"

	codec "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/codec"
	index "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/index"

	bolt "go.etcd.io/bbolt"
)

type Options struct {
	codec       codec.MarshalUnmarshaler
	boltMode    os.FileMode
	boltOptions *bolt.Options
	batchMode   bool
	rootBucket  []string
	path        string
	bolt        *bolt.DB
}

func BoltOptions(mode os.FileMode, options *bolt.Options) func(*Options) error {
	return func(opts *Options) error {
		opts.boltMode = mode
		opts.boltOptions = options
		return nil
	}
}

func Codec(c codec.MarshalUnmarshaler) func(*Options) error {
	return func(opts *Options) error {
		opts.codec = c
		return nil
	}
}

func Batch() func(*Options) error {
	return func(opts *Options) error {
		opts.batchMode = true
		return nil
	}
}

func Root(root ...string) func(*Options) error {
	return func(opts *Options) error {
		opts.rootBucket = root
		return nil
	}
}

func UseDB(b *bolt.DB) func(*Options) error {
	return func(opts *Options) error {
		opts.path = b.Path()
		opts.bolt = b
		return nil
	}
}

func Limit(limit int) func(*index.Options) {
	return func(opts *index.Options) {
		opts.Limit = limit
	}
}

func Skip(offset int) func(*index.Options) {
	return func(opts *index.Options) {
		opts.Skip = offset
	}
}

func Reverse() func(*index.Options) {
	return func(opts *index.Options) {
		opts.Reverse = true
	}
}
