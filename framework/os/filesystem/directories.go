package filesystem

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

func (self ApplicationDirectories) ParseApplicationDirectories() {
	var err error
	if self.Working, err = os.Getwd(); err != nil {
		panic(fmt.Sprintf("[fatal error] failed to determine working directory:", err))
	}
	self.Temporary = os.TempDir()
}

func (self *ApplicationDirectories) ParseUserDirectories() {
	var err error
	self.UserHome = os.Getenv("HOME")
	// TODO: Why is this undefined?
	// REF: https://golang.org/src/os/file.go
	//self.UserHomeDirectory, err = os.UserHomeDir()
	//if err != nil {
	//	panic(fmt.Sprintf("[fatal error] failed to determine user home:", err))
	//}
	if self.UserCache, err = os.UserCacheDir(); err != nil {
		panic(fmt.Sprintf("[fatal error] failed to determine user cache:", err))
	}
	if self.UserConfig = self.UserHome + "/.config/starship"; err != nil {
		panic(fmt.Sprintf("[fatal error] failed to determine user config path:", err))
	}
	if _, err := os.Stat(self.UserConfig); os.IsNotExist(err) {
		os.Mkdir(self.UserConfig, 0770)
	}
	if self.UserData = self.UserHome + "/.local/share/starship/"; err != nil {
		panic(fmt.Sprintf("[fatal error] failed to determine user data path:", err))
	}
	if _, err := os.Stat(self.UserData); os.IsNotExist(err) {
		os.Mkdir(self.UserData, 0770)
	}
}
