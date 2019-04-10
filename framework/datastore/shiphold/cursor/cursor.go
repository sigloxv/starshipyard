package cursor

import (
	"bytes"

	bolt "go.etcd.io/bbolt"
)

type Cursor struct {
	C       *bolt.Cursor
	Reverse bool
}

func (self *Cursor) First() ([]byte, []byte) {
	if self.Reverse {
		return self.C.Last()
	}
	return self.C.First()
}

func (self *Cursor) Next() ([]byte, []byte) {
	if self.Reverse {
		return self.C.Prev()
	}
	return self.C.Next()
}

type RangeCursor struct {
	C         *bolt.Cursor
	Reverse   bool
	Min       []byte
	Max       []byte
	CompareFn func([]byte, []byte) int
}

func (self *RangeCursor) First() ([]byte, []byte) {
	if self.Reverse {
		k, v := self.C.Seek(self.Max)
		// NOTE: If Seek doesn't find a key it goes to the next.
		// If so, we need to get the previous one to avoid
		// including bigger values. #218
		if !bytes.HasPrefix(k, self.Max) && k != nil {
			k, v = self.C.Prev()
		}
		return k, v
	}
	return self.C.Seek(self.Min)
}

func (self *RangeCursor) Next() ([]byte, []byte) {
	if self.Reverse {
		return self.C.Prev()
	}
	return self.C.Next()
}

func (self *RangeCursor) Continue(val []byte) bool {
	if self.Reverse {
		return val != nil && self.CompareFn(val, self.Min) >= 0
	}
	return val != nil && self.CompareFn(val, self.Max) <= 0
}

type PrefixCursor struct {
	C       *bolt.Cursor
	Reverse bool
	Prefix  []byte
}

func (self *PrefixCursor) First() ([]byte, []byte) {
	var k, v []byte
	// TODO: What is this, why would we iterate over a set to do nothing with it?
	for k, v = self.C.First(); k != nil && !bytes.HasPrefix(k, self.Prefix); k, v = self.C.Next() {
	}
	if k == nil {
		return nil, nil
	}
	if self.Reverse {
		kc, vc := k, v
		// TODO: And this?
		for ; kc != nil && bytes.HasPrefix(kc, self.Prefix); kc, vc = self.C.Next() {
			k, v = kc, vc
		}
		if kc != nil {
			k, v = self.C.Prev()
		}
	}
	return k, v
}

func (self *PrefixCursor) Next() ([]byte, []byte) {
	if self.Reverse {
		return self.C.Prev()
	}
	return self.C.Next()
}

func (self *PrefixCursor) Continue(val []byte) bool {
	return val != nil && bytes.HasPrefix(val, self.Prefix)
}
