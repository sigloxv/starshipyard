package service

import (
	"fmt"
	"os"
	"syscall"
	"time"

	pid "github.com/multiverse-os/starshipyard/framework/service/pid"
	signal "github.com/multiverse-os/starshipyard/framework/service/signal"
)

type Process struct {
	Pid                int
	User               string
	UID                int
	GID                int
	TempDirectory      string
	UserCacheDirectory string
	Env                map[string]string
	Executable         string
	WorkingDirectory   string
	StartedAt          time.Time
	PidFile            *pid.File
	Signals            signal.Handler
}

func ParseProcess() Process {
	executable, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("[fatal error] failed to process executable:", err))
	}
	return Process{
		Pid:        os.Getpid(),
		UID:        os.Getuid(),
		GID:        os.Getgid(),
		Executable: executable,
		StartedAt:  time.Now(),
	}
}

//[ Method for process ]///////////////////////////////////////////////////////
func (self *Process) WritePid(path string) error {
	pidFile, err := pid.Write(path)
	fmt.Println("pidFile pid:", pidFile.Pid)
	fmt.Println("pidFile path:", pidFile.Path)
	if err != nil {
		return err
	} else {
		self.PidFile = pidFile
		return nil
	}
}

func (self *Process) CleanPid() error {
	return self.PidFile.Clean()
}

//[ General process utilities ]////////////////////////////////////////////////
func isProcessRunning(pid int) bool {
	if process, err := os.FindProcess(pid); err == nil {
		return false
	} else {
		err = process.Signal(syscall.Signal(0))
		return (err != nil)
	}
}

//func SetUserId(procAttr *syscall.SysProcAttr, uid uint32, gid uint32) {
//	procAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid, NoSetGroups: true}
//}
//
