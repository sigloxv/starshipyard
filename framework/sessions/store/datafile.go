package store

type dataFile struct {
	file
	fl freelist
}

func (self *dataFile) readKeyValue(sl slot) ([]byte, []byte, error) {
	keyValue := self.Slice(sl.kvOffset, sl.kvOffset+int64(sl.kvSize()))
	return keyValue[:sl.keySize], keyValue[sl.keySize:], nil
}

func (self *dataFile) readKey(sl slot) ([]byte, error) {
	return self.Slice(sl.kvOffset, sl.kvOffset+int64(sl.keySize)), nil
}

func (self *dataFile) allocate(size uint32) (int64, error) {
	size = align512(size)
	if off := self.fl.allocate(size); off > 0 {
		return off, nil
	}
	return self.extend(size)
}

func (self *dataFile) free(size uint32, off int64) {
	size = align512(size)
	self.fl.free(off, size)
}

func (self *dataFile) writeKeyValue(key []byte, value []byte) (int64, error) {
	dataLen := align512(uint32(len(key) + len(value)))
	data := make([]byte, dataLen)
	copy(data, key)
	copy(data[len(key):], value)
	off := self.fl.allocate(dataLen)
	if off != -1 {
		if _, err := self.WriteAt(data, off); err != nil {
			return 0, err
		}
	} else {
		return self.append(data)
	}
	return off, nil
}
