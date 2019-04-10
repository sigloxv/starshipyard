package shiphold

import (
	"bytes"

	bolt "go.etcd.io/bbolt"
)

type BucketScanner interface {
	PrefixScan(prefix string) []Node
	RangeScan(min, max string) []Node
}

func (self *node) PrefixScan(prefix string) []Node {
	if self.tx != nil {
		return self.prefixScan(self.tx, prefix)
	}
	var nodes []Node
	self.readTx(func(tx *bolt.Tx) error {
		nodes = self.prefixScan(tx, prefix)
		return nil
	})
	return nodes
}

func (self *node) prefixScan(tx *bolt.Tx, prefix string) []Node {
	var (
		prefixBytes = []byte(prefix)
		nodes       []Node
		c           = self.cursor(tx)
	)
	for k, v := c.Seek(prefixBytes); k != nil && bytes.HasPrefix(k, prefixBytes); k, v = c.Next() {
		if v != nil {
			continue
		}
		nodes = append(nodes, self.From(string(k)))
	}
	return nodes
}

func (self *node) RangeScan(min, max string) []Node {
	if self.tx != nil {
		return self.rangeScan(self.tx, min, max)
	}
	var nodes []Node
	self.readTx(func(tx *bolt.Tx) error {
		nodes = self.rangeScan(tx, min, max)
		return nil
	})
	return nodes
}

func (self *node) rangeScan(tx *bolt.Tx, min, max string) []Node {
	var (
		minBytes = []byte(min)
		maxBytes = []byte(max)
		nodes    []Node
		c        = self.cursor(tx)
	)
	for k, v := c.Seek(minBytes); k != nil && bytes.Compare(k, maxBytes) <= 0; k, v = c.Next() {
		if v != nil {
			continue
		}
		nodes = append(nodes, self.From(string(k)))
	}
	return nodes
}

func (self *node) cursor(tx *bolt.Tx) *bolt.Cursor {
	var c *bolt.Cursor
	if len(self.rootBucket) > 0 {
		c = self.GetBucket(tx).Cursor()
	} else {
		c = tx.Cursor()
	}
	return c
}
