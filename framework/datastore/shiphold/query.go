package shiphold

import (
	cursor "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/cursor"
	q "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/q"

	bolt "go.etcd.io/bbolt"
)

func (self *node) Select(matchers ...q.Matcher) Query {
	tree := q.And(matchers...)
	return newQuery(self, tree)
}

type Query interface {
	Skip(int) Query
	Limit(int) Query
	OrderBy(...string) Query
	Reverse() Query
	Bucket(string) Query
	Find(interface{}) error
	First(interface{}) error
	Delete(interface{}) error
	Count(interface{}) (int, error)
	Raw() ([][]byte, error)
	RawEach(func([]byte, []byte) error) error
	Each(interface{}, func(interface{}) error) error
}

func newQuery(n *node, tree q.Matcher) *query {
	return &query{
		skip:  0,
		limit: -1,
		node:  n,
		tree:  tree,
	}
}

type query struct {
	limit   int
	skip    int
	reverse bool
	tree    q.Matcher
	node    *node
	bucket  string
	orderBy []string
}

func (self *query) Skip(nb int) Query {
	self.skip = nb
	return self
}

func (self *query) Limit(nb int) Query {
	self.limit = nb
	return self
}

func (self *query) OrderBy(field ...string) Query {
	self.orderBy = field
	return self
}

func (self *query) Reverse() Query {
	self.reverse = true
	return self
}

func (self *query) Bucket(bucketName string) Query {
	self.bucket = bucketName
	return self
}

func (self *query) Find(to interface{}) error {
	sink, err := newListSink(self.node, to)
	if err != nil {
		return err
	}
	return self.runQuery(sink)
}

func (self *query) First(to interface{}) error {
	sink, err := newFirstSink(self.node, to)
	if err != nil {
		return err
	}
	self.limit = 1
	return self.runQuery(sink)
}

func (self *query) Delete(kind interface{}) error {
	sink, err := newDeleteSink(self.node, kind)
	if err != nil {
		return err
	}

	return self.runQuery(sink)
}

func (self *query) Count(kind interface{}) (int, error) {
	sink, err := newCountSink(self.node, kind)
	if err != nil {
		return 0, err
	}

	err = self.runQuery(sink)
	if err != nil {
		return 0, err
	}

	return sink.counter, nil
}

func (self *query) Raw() ([][]byte, error) {
	sink := newRawSink()
	if err := self.runQuery(sink); err != nil {
		return nil, err
	}
	return sink.results, nil
}

func (self *query) RawEach(fn func([]byte, []byte) error) error {
	sink := newRawSink()
	sink.execFn = fn
	return self.runQuery(sink)
}

func (self *query) Each(kind interface{}, fn func(interface{}) error) error {
	sink, err := newEachSink(kind)
	if err != nil {
		return err
	}
	sink.execFn = fn
	return self.runQuery(sink)
}

func (self *query) runQuery(sink sink) error {
	if self.node.tx != nil {
		return self.query(self.node.tx, sink)
	}
	if sink.readOnly() {
		return self.node.s.Bolt.View(func(tx *bolt.Tx) error {
			return self.query(tx, sink)
		})
	}
	return self.node.s.Bolt.Update(func(tx *bolt.Tx) error {
		return self.query(tx, sink)
	})
}

func (self *query) query(tx *bolt.Tx, sink sink) error {
	bucketName := self.bucket
	if bucketName == "" {
		bucketName = sink.bucketName()
	}
	bucket := self.node.GetBucket(tx, bucketName)

	if self.limit == 0 {
		return sink.flush()
	}
	sorter := newSorter(self.node, sink)
	sorter.orderBy = self.orderBy
	sorter.reverse = self.reverse
	sorter.skip = self.skip
	sorter.limit = self.limit
	if bucket != nil {
		c := cursor.Cursor{C: bucket.Cursor(), Reverse: self.reverse}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v == nil {
				continue
			}
			stop, err := sorter.filter(self.tree, bucket, k, v)
			if err != nil {
				return err
			}
			if stop {
				break
			}
		}
	}
	return sorter.flush()
}
