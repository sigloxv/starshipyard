package framework

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/unix"
)

func lockFcntl(name string) (io.Closer, error) {
	fi, err := os.Stat(name)
	if err == nil && fi.Size() > 0 {
		return nil, fmt.Errorf("can't Lock file %q: has non-zero size", name)
	}

	f, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("Lock Create of %s failed: %v", name, err)
	}

	err = unix.FcntlFlock(f.Fd(), unix.F_SETLK, &unix.Flock_t{
		Type:   unix.F_WRLCK,
		Whence: int16(os.SEEK_SET),
		Start:  0,
		Len:    0, // 0 means to lock the entire file.
		Pid:    0, // only used by F_GETLK
	})

	if err != nil {
		f.Close()
		return nil, fmt.Errorf("Lock FcntlFlock of %s failed: %v", name, err)
	}
	return &unlocker{f: f, abs: name}, nil
}
