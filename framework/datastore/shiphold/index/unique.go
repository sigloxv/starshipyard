package index

import (
	"bytes"

	cursor "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/cursor"

	bolt "go.etcd.io/bbolt"
)

func NewUniqueIndex(parent *bolt.Bucket, indexName []byte) (*UniqueIndex, error) {
	var err error
	b := parent.Bucket(indexName)
	if b == nil {
		if !parent.Writable() {
			return nil, errNotFound
		}
		b, err = parent.CreateBucket(indexName)
		if err != nil {
			return nil, err
		}
	}
	return &UniqueIndex{
		IndexBucket: b,
		Parent:      parent,
	}, nil
}

type UniqueIndex struct {
	Parent      *bolt.Bucket
	IndexBucket *bolt.Bucket
}

func (self *UniqueIndex) Add(value []byte, targetID []byte) error {
	if value == nil || len(value) == 0 {
		return errNilParam
	}
	if targetID == nil || len(targetID) == 0 {
		return errNilParam
	}

	exists := self.IndexBucket.Get(value)
	if exists != nil {
		if bytes.Equal(exists, targetID) {
			return nil
		}
		return errAlreadyExists
	}

	return self.IndexBucket.Put(value, targetID)
}

func (self *UniqueIndex) Remove(value []byte) error {
	return self.IndexBucket.Delete(value)
}

func (self *UniqueIndex) RemoveID(id []byte) error {
	c := self.IndexBucket.Cursor()

	for val, ident := c.First(); val != nil; val, ident = c.Next() {
		if bytes.Equal(ident, id) {
			return self.Remove(val)
		}
	}
	return nil
}

func (self *UniqueIndex) Get(value []byte) []byte {
	return self.IndexBucket.Get(value)
}

func (self *UniqueIndex) All(value []byte, opts *Options) ([][]byte, error) {
	id := self.IndexBucket.Get(value)
	if id != nil {
		return [][]byte{id}, nil
	}
	return nil, nil
}

func (self *UniqueIndex) AllRecords(opts *Options) ([][]byte, error) {
	var list [][]byte

	c := cursor.Cursor{C: self.IndexBucket.Cursor(), Reverse: opts != nil && opts.Reverse}

	for val, ident := c.First(); val != nil; val, ident = c.Next() {
		if opts != nil && opts.Skip > 0 {
			opts.Skip--
			continue
		}
		if opts != nil && opts.Limit == 0 {
			break
		}
		if opts != nil && opts.Limit > 0 {
			opts.Limit--
		}
		list = append(list, ident)
	}
	return list, nil
}

func (self *UniqueIndex) Range(min []byte, max []byte, opts *Options) ([][]byte, error) {
	var list [][]byte

	c := cursor.RangeCursor{
		C:       self.IndexBucket.Cursor(),
		Reverse: opts != nil && opts.Reverse,
		Min:     min,
		Max:     max,
		CompareFn: func(val, limit []byte) int {
			return bytes.Compare(val, limit)
		},
	}

	for val, ident := c.First(); val != nil && c.Continue(val); val, ident = c.Next() {
		if opts != nil && opts.Skip > 0 {
			opts.Skip--
			continue
		}
		if opts != nil && opts.Limit == 0 {
			break
		}
		if opts != nil && opts.Limit > 0 {
			opts.Limit--
		}
		list = append(list, ident)
	}
	return list, nil
}

func (self *UniqueIndex) Prefix(prefix []byte, opts *Options) ([][]byte, error) {
	var list [][]byte

	c := cursor.PrefixCursor{
		C:       self.IndexBucket.Cursor(),
		Reverse: opts != nil && opts.Reverse,
		Prefix:  prefix,
	}

	for val, ident := c.First(); val != nil && c.Continue(val); val, ident = c.Next() {
		if opts != nil && opts.Skip > 0 {
			opts.Skip--
			continue
		}
		if opts != nil && opts.Limit == 0 {
			break
		}
		if opts != nil && opts.Limit > 0 {
			opts.Limit--
		}
		list = append(list, ident)
	}
	return list, nil
}

func (self *UniqueIndex) first() []byte {
	c := self.IndexBucket.Cursor()
	for val, ident := c.First(); val != nil; val, ident = c.Next() {
		return ident
	}
	return nil
}
