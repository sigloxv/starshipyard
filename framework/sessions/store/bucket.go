package store

import (
	"encoding/binary"

	fs "github.com/multiverse-os/starshipyard/framework/sessions/store/fs"
)

type slot struct {
	hash      uint32
	keySize   uint16
	valueSize uint32
	kvOffset  int64
}

func (self slot) kvSize() uint32 {
	return uint32(self.keySize) + self.valueSize
}

type bucket struct {
	slots [slotsPerBucket]slot
	next  int64
}

type bucketHandle struct {
	bucket
	file   fs.MmapFile
	offset int64
}

const (
	bucketSize uint32 = 512
)

func align512(n uint32) uint32 {
	return (n + 511) &^ 511
}

func (self bucket) MarshalBinary() ([]byte, error) {
	buf := make([]byte, bucketSize)
	data := buf
	for i := 0; i < slotsPerBucket; i++ {
		sl := self.slots[i]
		binary.LittleEndian.PutUint32(buf[:4], sl.hash)
		binary.LittleEndian.PutUint16(buf[4:6], sl.keySize)
		binary.LittleEndian.PutUint32(buf[6:10], sl.valueSize)
		binary.LittleEndian.PutUint64(buf[10:18], uint64(sl.kvOffset))
		buf = buf[18:]
	}
	binary.LittleEndian.PutUint64(buf[:8], uint64(self.next))
	return data, nil
}

func (self *bucket) UnmarshalBinary(data []byte) error {
	for i := 0; i < slotsPerBucket; i++ {
		_ = data[18] // bounds check hint to compiler; see golang.org/issue/14808
		self.slots[i].hash = binary.LittleEndian.Uint32(data[:4])
		self.slots[i].keySize = binary.LittleEndian.Uint16(data[4:6])
		self.slots[i].valueSize = binary.LittleEndian.Uint32(data[6:10])
		self.slots[i].kvOffset = int64(binary.LittleEndian.Uint64(data[10:18]))
		data = data[18:]
	}
	self.next = int64(binary.LittleEndian.Uint64(data[:8]))
	return nil
}

func (self *bucket) del(slotIdx int) {
	i := slotIdx
	for ; i < slotsPerBucket-1; i++ {
		self.slots[i] = self.slots[i+1]
	}
	self.slots[i] = slot{}
}

func (self *bucketHandle) read() error {
	buf := self.file.Slice(self.offset, self.offset+int64(bucketSize))
	return self.UnmarshalBinary(buf)
}

func (self *bucketHandle) write() error {
	buf, err := self.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = self.file.WriteAt(buf, self.offset)
	return err
}

type slotWriter struct {
	bucket      *bucketHandle
	slotIdx     int
	prevBuckets []*bucketHandle
}

func (self *slotWriter) insert(sl slot, db *DB) error {
	if self.slotIdx == slotsPerBucket {
		nextBucket, err := db.createOverflowBucket()
		if err != nil {
			return err
		}
		self.bucket.next = nextBucket.offset
		self.prevBuckets = append(self.prevBuckets, self.bucket)
		self.bucket = nextBucket
		self.slotIdx = 0
	}
	self.bucket.slots[self.slotIdx] = sl
	self.slotIdx++
	return nil
}

func (self *slotWriter) write() error {
	for i := len(self.prevBuckets) - 1; i >= 0; i-- {
		if err := self.prevBuckets[i].write(); err != nil {
			return err
		}
	}
	return self.bucket.write()
}
