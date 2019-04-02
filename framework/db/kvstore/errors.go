package kvstore

import "errors"

var (
	ErrBucketNotFound = errors.New("Bucket not found")
	ErrKeyNotFound    = errors.New("Key not found")
	ErrDoesNotExist   = errors.New("Does not exist")
	ErrFoundIt        = errors.New("Found it")
	ErrExistsInSet    = errors.New("Element already exists in set")
	ErrInvalidID      = errors.New("Element ID can not contain \":\"")
)
