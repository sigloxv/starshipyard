package shiphold

import (
	bolt "go.etcd.io/bbolt"
)

func (self *node) CreateBucketIfNotExists(tx *bolt.Tx, bucket string) (b *bolt.Bucket, err error) {
	bucketNames := append(self.rootBucket, bucket)
	for _, bucketName := range bucketNames {
		if b != nil {
			if b, err = b.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
				return nil, err
			}

		} else {
			if b, err = tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
				return nil, err
			}
		}
	}
	return b, nil
}

func (self *node) GetBucket(tx *bolt.Tx, children ...string) (b *bolt.Bucket) {
	bucketNames := append(self.rootBucket, children...)
	for _, bucketName := range bucketNames {
		if b != nil {
			if b = b.Bucket([]byte(bucketName)); b == nil {
				return nil
			}
		} else {
			if b = tx.Bucket([]byte(bucketName)); b == nil {
				return nil
			}
		}
	}
	return b
}
