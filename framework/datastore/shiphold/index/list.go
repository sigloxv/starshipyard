package index

import (
	"bytes"

	cursor "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/cursor"

	bolt "go.etcd.io/bbolt"
)

func NewListIndex(parent *bolt.Bucket, indexName []byte) (*ListIndex, error) {
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

	ids, err := NewUniqueIndex(b, []byte("shiphold__ids"))
	if err != nil {
		return nil, err
	}

	return &ListIndex{
		IndexBucket: b,
		Parent:      parent,
		IDs:         ids,
	}, nil
}

type ListIndex struct {
	Parent      *bolt.Bucket
	IndexBucket *bolt.Bucket
	IDs         *UniqueIndex
}

func (self *ListIndex) Add(newValue []byte, targetID []byte) error {
	if newValue == nil || len(newValue) == 0 {
		return errNilParam
	}
	if targetID == nil || len(targetID) == 0 {
		return errNilParam
	}

	key := self.IDs.Get(targetID)
	if key != nil {
		err := self.IndexBucket.Delete(key)
		if err != nil {
			return err
		}

		err = self.IDs.Remove(targetID)
		if err != nil {
			return err
		}

		key = key[:0]
	}

	key = append(key, newValue...)
	key = append(key, '_')
	key = append(key, '_')
	key = append(key, targetID...)

	if err := self.IDs.Add(targetID, key); err != nil {
		return err
	}
	return self.IndexBucket.Put(key, targetID)
}

func (self *ListIndex) Remove(value []byte) error {
	var err error
	var keys [][]byte

	c := self.IndexBucket.Cursor()
	prefix := generatePrefix(value)

	for k, _ := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, _ = c.Next() {
		keys = append(keys, k)
	}
	for _, k := range keys {
		err = self.IndexBucket.Delete(k)
		if err != nil {
			return err
		}
	}
	return self.IDs.RemoveID(value)
}

func (self *ListIndex) RemoveID(targetID []byte) error {
	value := self.IDs.Get(targetID)
	if value == nil {
		return nil
	}

	err := self.IndexBucket.Delete(value)
	if err != nil {
		return err
	}

	return self.IDs.Remove(targetID)
}

func (self *ListIndex) Get(value []byte) []byte {
	c := self.IndexBucket.Cursor()
	prefix := generatePrefix(value)

	for k, id := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, id = c.Next() {
		return id
	}

	return nil
}

func (self *ListIndex) All(value []byte, opts *Options) ([][]byte, error) {
	var list [][]byte
	c := self.IndexBucket.Cursor()
	cur := cursor.Cursor{C: c, Reverse: opts != nil && opts.Reverse}

	prefix := generatePrefix(value)

	k, id := c.Seek(prefix)
	if cur.Reverse {
		var count int
		kc := k
		idc := id
		for ; kc != nil && bytes.HasPrefix(kc, prefix); kc, idc = c.Next() {
			count++
			k, id = kc, idc
		}
		if kc != nil {
			k, id = c.Prev()
		}
		list = make([][]byte, 0, count)
	}

	for ; bytes.HasPrefix(k, prefix); k, id = cur.Next() {
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

		list = append(list, id)
	}

	return list, nil
}

// TODO: Does this actualy have any condition where it returns an error?
func (self *ListIndex) AllRecords(opts *Options) ([][]byte, error) {
	var list [][]byte

	c := cursor.Cursor{C: self.IndexBucket.Cursor(), Reverse: opts != nil && opts.Reverse}

	for k, id := c.First(); k != nil; k, id = c.Next() {
		if id == nil || bytes.Equal(k, []byte("shiphold__ids")) {
			continue
		}
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
		list = append(list, id)
	}
	return list, nil
}

func (self *ListIndex) Range(min []byte, max []byte, opts *Options) ([][]byte, error) {
	var list [][]byte

	c := cursor.RangeCursor{
		C:       self.IndexBucket.Cursor(),
		Reverse: opts != nil && opts.Reverse,
		Min:     min,
		Max:     max,
		CompareFn: func(val, limit []byte) int {
			pos := bytes.LastIndex(val, []byte("__"))
			return bytes.Compare(val[:pos], limit)
		},
	}

	for k, id := c.First(); c.Continue(k); k, id = c.Next() {
		if id == nil || bytes.Equal(k, []byte("shiphold__ids")) {
			continue
		}

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

		list = append(list, id)
	}

	return list, nil
}

func (self *ListIndex) Prefix(prefix []byte, opts *Options) ([][]byte, error) {
	var list [][]byte

	c := cursor.PrefixCursor{
		C:       self.IndexBucket.Cursor(),
		Reverse: opts != nil && opts.Reverse,
		Prefix:  prefix,
	}

	for k, id := c.First(); k != nil && c.Continue(k); k, id = c.Next() {
		if id == nil || bytes.Equal(k, []byte("shiphold__ids")) {
			continue
		}

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

		list = append(list, id)
	}
	return list, nil
}

func generatePrefix(value []byte) []byte {
	prefix := make([]byte, len(value)+2)
	var i int
	for i = range value {
		prefix[i] = value[i]
	}
	prefix[i+1] = '_'
	prefix[i+2] = '_'
	return prefix
}
