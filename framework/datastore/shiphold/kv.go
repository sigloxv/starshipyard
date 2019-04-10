package shiphold

import (
	"reflect"

	bolt "go.etcd.io/bbolt"
)

type KeyValueStore interface {
	Get(bucketName string, key interface{}, to interface{}) error
	Set(bucketName string, key interface{}, value interface{}) error
	Delete(bucketName string, key interface{}) error
	GetBytes(bucketName string, key interface{}) ([]byte, error)
	SetBytes(bucketName string, key interface{}, value []byte) error
	KeyExists(bucketName string, key interface{}) (bool, error)
}

func (self *node) GetBytes(bucketName string, key interface{}) ([]byte, error) {
	id, err := toBytes(key, self.codec)
	if err != nil {
		return nil, err
	}
	var val []byte
	return val, self.readTx(func(tx *bolt.Tx) error {
		raw, err := self.getBytes(tx, bucketName, id)
		if err != nil {
			return err
		}

		val = make([]byte, len(raw))
		copy(val, raw)
		return nil
	})
}

// GetBytes gets a raw value from a bucket.
func (self *node) getBytes(tx *bolt.Tx, bucketName string, id []byte) ([]byte, error) {
	bucket := self.GetBucket(tx, bucketName)
	if bucket == nil {
		return nil, errNotFound
	}

	raw := bucket.Get(id)
	if raw == nil {
		return nil, errNotFound
	}

	return raw, nil
}

// SetBytes sets a raw value into a bucket.
func (self *node) SetBytes(bucketName string, key interface{}, value []byte) error {
	if key == nil {
		return errNilParam
	}

	id, err := toBytes(key, self.codec)
	if err != nil {
		return err
	}

	return self.readWriteTx(func(tx *bolt.Tx) error {
		return self.setBytes(tx, bucketName, id, value)
	})
}

func (self *node) setBytes(tx *bolt.Tx, bucketName string, id, data []byte) error {
	bucket, err := self.CreateBucketIfNotExists(tx, bucketName)
	if err != nil {
		return err
	}
	if _, err = newMeta(bucket, self); err != nil {
		return err
	}
	return bucket.Put(id, data)
}

func (self *node) Get(bucketName string, key interface{}, to interface{}) error {
	ref := reflect.ValueOf(to)

	if !ref.IsValid() || ref.Kind() != reflect.Ptr {
		return errPtrNeeded
	}
	id, err := toBytes(key, self.codec)
	if err != nil {
		return err
	}
	return self.readTx(func(tx *bolt.Tx) error {
		raw, err := self.getBytes(tx, bucketName, id)
		if err != nil {
			return err
		}
		return self.codec.Unmarshal(raw, to)
	})
}

func (self *node) Set(bucketName string, key interface{}, value interface{}) error {
	var data []byte
	var err error
	if value != nil {
		data, err = self.codec.Marshal(value)
		if err != nil {
			return err
		}
	}
	return self.SetBytes(bucketName, key, data)
}

func (self *node) Delete(bucketName string, key interface{}) error {
	id, err := toBytes(key, self.codec)
	if err != nil {
		return err
	}
	return self.readWriteTx(func(tx *bolt.Tx) error {
		return self.delete(tx, bucketName, id)
	})
}

func (self *node) delete(tx *bolt.Tx, bucketName string, id []byte) error {
	bucket := self.GetBucket(tx, bucketName)
	if bucket == nil {
		return errNotFound
	}
	return bucket.Delete(id)
}

func (self *node) KeyExists(bucketName string, key interface{}) (bool, error) {
	id, err := toBytes(key, self.codec)
	if err != nil {
		return false, err
	}
	var exists bool
	return exists, self.readTx(func(tx *bolt.Tx) error {
		bucket := self.GetBucket(tx, bucketName)
		if bucket == nil {
			return errNotFound
		}
		v := bucket.Get(id)
		if v != nil {
			exists = true
		}
		return nil
	})
}
