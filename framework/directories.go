package framework

import (
	"fmt"
	"os"
)

func (self *Application) ParseApplicationDirectories() {
	var err error
	self.WorkingDirectory, err = os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("[fatal error] failed to determine working directory:", err))
	}
	self.TemporaryDirectory = os.TempDir()
	if err != nil {
		panic(fmt.Sprintf("[fatal error] failed to obtain temporary directory:", err))
	}
}

func (self *Application) ParseUserDirectories() {
	var err error
	self.UserHomeDirectory = os.Getenv("HOME")
	// TODO: Why is this undefined?
	// REF: https://golang.org/src/os/file.go
	//self.UserHomeDirectory, err = os.UserHomeDir()
	//if err != nil {
	//	panic(fmt.Sprintf("[fatal error] failed to determine user home:", err))
	//}
	self.UserCacheDirectory, err = os.UserCacheDir()
	if err != nil {
		panic(fmt.Sprintf("[fatal error] failed to determine user cache:", err))
	}

	self.UserConfigDirectory = self.UserHomeDirectory + "/.config/starship"
	if err != nil {
		panic(fmt.Sprintf("[fatal error] failed to determine user config path:", err))
	}
	if _, err := os.Stat(self.UserConfigDirectory); os.IsNotExist(err) {
		os.Mkdir(self.UserConfigDirectory, 0770)
	}

	self.UserDataDirectory = self.UserHomeDirectory + "/.local/share/starship/"
	if err != nil {
		panic(fmt.Sprintf("[fatal error] failed to determine user data path:", err))
	}
	if _, err := os.Stat(self.UserDataDirectory); os.IsNotExist(err) {
		os.Mkdir(self.UserDataDirectory, 0770)
	}
}
