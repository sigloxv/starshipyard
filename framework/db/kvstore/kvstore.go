package kvstore

import (
	"encoding/binary"
	"errors"
	"strconv"
	"time"

	bolt "github.com/etcd-io/bbolt"
	xid "github.com/rs/xid"
)

type (
	Database bolt.DB

	boltBucket struct {
		db   *Database // the Bolt database
		name []byte    // the bucket name
	}

	KeyValue boltBucket
)

type KVStore struct {
	Store       *Database
	Collections map[string]*KeyValue
}

type Entry struct {
	ID    xid.ID
	Key   string
	Bytes []byte
}

func New(filename string) (*Database, error) {
	// Use a timeout, in case the database file is already in use
	db, err := bolt.Open(filename, 0600, &bolt.Options{Timeout: 15 * time.Second})
	if err != nil {
		return nil, err
	}
	return (*Database)(db), nil
}

func (db *Database) Close() {
	(*bolt.DB)(db).Close()
}

func (self *Database) NewCollection(id string) (*KeyValue, error) {
	name := []byte(id)
	if err := (*bolt.DB)(self).Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(name); err != nil {
			return errors.New("Could not create bucket: " + err.Error())
		}
		return nil // Return from Update function
	}); err != nil {
		return nil, err
	}
	return &KeyValue{self, name}, nil
}

func (kv *KeyValue) Set(key string, value string) error {
	if kv.name == nil {
		return ErrDoesNotExist
	}
	return (*bolt.DB)(kv.db).Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(kv.name)
		if bucket == nil {
			return ErrBucketNotFound
		}
		return bucket.Put([]byte(key), []byte(value))
	})
}

func (kv *KeyValue) Get(key string) (value string, err error) {
	if kv.name == nil {
		return "", ErrDoesNotExist
	}
	err = (*bolt.DB)(kv.db).View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(kv.name)
		if bucket == nil {
			return ErrBucketNotFound
		}
		value := string(bucket.Get([]byte(key)))
		if value == "" {
			return ErrKeyNotFound
		} else {
			return nil // Return from View function
		}
	})
	return value, nil
}

func (kv *KeyValue) Delete(key string) error {
	if kv.name == nil {
		return ErrDoesNotExist
	}
	return (*bolt.DB)(kv.db).Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(kv.name)
		if bucket == nil {
			return ErrBucketNotFound
		}
		return bucket.Delete([]byte(key))
	})
}

func (kv *KeyValue) Increment(key string) (value int, err error) {
	if kv.name == nil {
		kv.name = []byte(key)
	}
	value = 0
	return value, (*bolt.DB)(kv.db).Update(func(tx *bolt.Tx) error {
		// The numeric value
		// Get the string value
		bucket := tx.Bucket(kv.name)
		if bucket == nil {
			// Create the bucket if it does not already exist
			bucket, err = tx.CreateBucketIfNotExists(kv.name)
			if err != nil {
				return errors.New("Could not create bucket: " + err.Error())
			}
		} else {
			value, err = strconv.Atoi(string(bucket.Get([]byte(key))))
		}
		value++
		err = bucket.Put([]byte(key), []byte(strconv.Itoa(value)))
		return err
	})
}

func (kv *KeyValue) Remove() error {
	err := (*bolt.DB)(kv.db).Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(kv.name))
	})
	// Mark as removed by setting the name to nil
	kv.name = nil
	return err
}

func (kv *KeyValue) Clear() error {
	if kv.name == nil {
		return ErrDoesNotExist
	}
	return (*bolt.DB)(kv.db).Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(kv.name)
		if bucket == nil {
			return ErrBucketNotFound
		}
		return bucket.ForEach(func(key, _ []byte) error {
			return bucket.Delete(key)
		})
	})
}

func byteID(x uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, x)
	return b
}
