// +build !linux !arm64

package daemon

import (
	"syscall"
)

func syscallDup(oldfd int, newfd int) (err error) {
	return syscall.Dup2(oldfd, newfd)
}
