package fs

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	initialMmapSize = 1024 << 20
)

type osfile struct {
	*os.File
	data     []byte
	mmapSize int64
}

type osfs struct{}

// OS is a file system backed by the os package.
var OS = &osfs{}

func (fs *osfs) OpenFile(name string, flag int, perm os.FileMode) (MmapFile, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	mf := &osfile{f, nil, 0}
	if stat.Size() > 0 {
		if err := mf.Mmap(stat.Size()); err != nil {
			return nil, err
		}
	}
	return mf, err
}

func (fs *osfs) CreateLockFile(name string, perm os.FileMode) (LockFile, bool, error) {
	return createLockFile(name, perm)
}

func (fs *osfs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (fs *osfs) Remove(name string) error {
	return os.Remove(name)
}

type oslockfile struct {
	File
	path string
}

func (f *oslockfile) Unlock() error {
	if err := os.Remove(f.path); err != nil {
		return err
	}
	return f.Close()
}

func (f *osfile) Slice(start int64, end int64) []byte {
	return f.data[start:end]
}

func (f *osfile) Close() error {
	if err := munmap(f.data); err != nil {
		return nil
	}
	f.data = nil
	return f.File.Close()
}

func (f *osfile) Mmap(fileSize int64) error {
	mmapSize := f.mmapSize

	if mmapSize >= fileSize {
		return nil
	}

	if mmapSize == 0 {
		mmapSize = initialMmapSize
		if mmapSize < fileSize {
			mmapSize = fileSize
		}
	} else {
		if err := munmap(f.data); err != nil {
			return err
		}
		mmapSize *= 2
	}

	data, mappedSize, err := mmap(f.File, fileSize, mmapSize)
	if err != nil {
		return err
	}

	madviceRandom(data)

	f.data = data
	f.mmapSize = mappedSize
	return nil
}

func mmap(f *os.File, fileSize int64, mmapSize int64) ([]byte, int64, error) {
	p, err := syscall.Mmap(int(f.Fd()), 0, int(mmapSize), syscall.PROT_READ, syscall.MAP_SHARED)
	return p, mmapSize, err
}

func munmap(data []byte) error {
	return syscall.Munmap(data)
}

func madviceRandom(data []byte) error {
	_, _, errno := syscall.Syscall(syscall.SYS_MADVISE, uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)), uintptr(syscall.MADV_RANDOM))
	if errno != 0 {
		return errno
	}
	return nil
}

func createLockFile(name string, perm os.FileMode) (LockFile, bool, error) {
	acquiredExisting := false
	if _, err := os.Stat(name); err == nil {
		acquiredExisting = true
	}
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		return nil, false, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		if err == syscall.EWOULDBLOCK {
			err = os.ErrExist
		}
		return nil, false, err
	}
	return &oslockfile{f, name}, acquiredExisting, nil
}
