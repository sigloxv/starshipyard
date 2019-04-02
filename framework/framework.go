package framework

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	scramble "github.com/multiverse-os/scramble-key"
	config "github.com/multiverse-os/starshipyard/framework/config"
	db "github.com/multiverse-os/starshipyard/framework/db"
	server "github.com/multiverse-os/starshipyard/framework/server"
	service "github.com/multiverse-os/starshipyard/framework/service"
)

// NOTE: Concept: we want to be able to run multiple applications in a given
// instance. This would likely be defined by a ruby-like script config that
// defines what domains go where, reverse and inverting proxy settings, etc
type Domain struct {
	Name        string
	Subdomains  []string
	Certificate string
}

// TODO: Starship Yard is meant to function as both the web application and
// include the firewall/reverse-proxy. This means we need to be able to register
// multiple end-points.
// TODO: Migrate to less pointer use, we should only be using pointers when
// absolutely necessary to do so, or when there is an obvious benefit. It should
// absolutely not be the default, just like globals.
// TODO: Need to persist and load scramble key
// TODO: Build a string function to provide a nice ouput with all necesary
// information
type Application struct {
	Name                string
	Process             service.Process
	ParentPid           int
	User                string
	UID                 int
	WorkingDirectory    string
	UserHomeDirectory   string
	UserCacheDirectory  string
	UserConfigDirectory string
	UserDataDirectory   string
	TemporaryDirectory  string
	Config              *config.Config
	ScrambleKey         scramble.Key
	KV                  *db.KV
	HTTPServer          *server.Server
}

// TODO: Want to use function chaning when initializing the server so that
// routes can be passed in
func Init(config *config.Config) *Application {
	rand.Seed(time.Now().UTC().UnixNano())
	service.WritePid(config.Pid)
	app := &Application{
		Config:      config,
		HTTPServer:  server.New(config),
		KV:          db.InitKV(),
		ScrambleKey: scramble.GenerateKey(),
		Process:     service.ParseProcess(),
	}

	app.Process.Signals = service.OnShutdownSignals(func(s os.Signal) {
		fmt.Println("[starship] received exit signal:", s)
		app.Stop()
	})

	app.Process.WritePid(app.Config.Pid)

	app.ParseApplicationDirectories()
	//app.ParseUserDirectories()
	app.Process = service.ParseProcess()
	app.KV.NewCollection("users")
	// TODO: Connect/Initialize/Load databases
	fmt.Println("application:", app)
	return app
}

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

func (self *Application) CleanUp() error {
	fmt.Println("[starship] attempting to exit cleanly...")
	fmt.Println("[starship] closing the general key/value store")
	self.KV.Store.Close()
	fmt.Println("[starship] cleaning the pid file")
	err := self.Process.CleanPid()
	if err != nil {
		return err
	}
	return nil
}
