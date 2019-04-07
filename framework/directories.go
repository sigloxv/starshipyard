package framework

import (
	"fmt"
	"os"
)

type ApplicationDirectories struct {
	Working    string
	Data       string
	Temporary  string
	UserHome   string
	UserCache  string
	UserConfig string
	UserData   string
}

func (self *Application) ParseApplicationDirectories() {
	var err error
	if self.Directories.Working, err = os.Getwd(); err != nil {
		panic(fmt.Sprintf("[fatal error] failed to determine working directory:", err))
	}
	if self.Directories.Temporary = os.TempDir(); err != nil {
		panic(fmt.Sprintf("[fatal error] failed to obtain temporary directory:", err))
	}
}

func (self *Application) ParseUserDirectories() {
	var err error
	self.Directories.UserHome = os.Getenv("HOME")
	// TODO: Why is this undefined?
	// REF: https://golang.org/src/os/file.go
	//self.UserHomeDirectory, err = os.UserHomeDir()
	//if err != nil {
	//	panic(fmt.Sprintf("[fatal error] failed to determine user home:", err))
	//}
	if self.Directories.UserCache, err = os.UserCacheDir(); err != nil {
		panic(fmt.Sprintf("[fatal error] failed to determine user cache:", err))
	}
	if self.Directories.UserConfig = self.Directories.UserHome + "/.config/starship"; err != nil {
		panic(fmt.Sprintf("[fatal error] failed to determine user config path:", err))
	}
	if _, err := os.Stat(self.Directories.UserConfig); os.IsNotExist(err) {
		os.Mkdir(self.Directories.UserConfig, 0770)
	}
	if self.Directories.UserData = self.Directories.UserHome + "/.local/share/starship/"; err != nil {
		panic(fmt.Sprintf("[fatal error] failed to determine user data path:", err))
	}
	if _, err := os.Stat(self.Directories.UserData); os.IsNotExist(err) {
		os.Mkdir(self.Directories.UserData, 0770)
	}
}
