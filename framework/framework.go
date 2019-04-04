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
	Name               string
	Process            *service.Process
	ParentPid          int
	User               string
	UID                int
	WorkingDirectory   string
	DataDirectory      string
	TemporaryDirectory string
	Config             *config.Config
	ScrambleKey        scramble.Key
	KV                 *db.KV
	HTTPServer         *server.Server

	UserHomeDirectory   string
	UserCacheDirectory  string
	UserConfigDirectory string
	UserDataDirectory   string
}

func Reseed() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func Init(config *config.Config) *Application {
	// TODO: Need to validate the values coming in, we need to ensure that values
	// that absolutey can not be nil like pid have defaults to fallback on
	rand.Seed(time.Now().UTC().UnixNano())

	if service.IsRootUser() {
		fmt.Println("[starship] running internet facing servers as root is very dangerous, run as an unpriviledged user")
		os.Exit(1)
		// TODO: Lets drop priviledges
	}

	wd, _ := os.Getwd()

	// NOTE: This is bare minimum validation and default fallbacks so that errors
	// are not thrown when setting up the application process signal handler, pid
	// control and other service functionality. An improved string/path validation
	// needs to be built ontop fo this basic functionality
	if len(config.Pid) == 0 {
		config.Pid = "tmp/pids/starship.pid"
	}
	if len(config.TemporaryDirectory) == 0 {
		config.TemporaryDirectory = "tmp"
	}
	if len(config.DataDirectory) == 0 {
		config.DataDirectory = "data"
	}
	if len(config.CacheDirectory) == 0 {
		config.CacheDirectory = "tmp/cache"
	}

	service.WritePid(config.Pid)

	app := &Application{
		Config:           config,
		WorkingDirectory: wd,
		HTTPServer:       server.New(config),
		KV:               db.InitKV(),
		ScrambleKey:      scramble.GenerateKey(),
		Process:          service.ParseProcess(),
	}

	app.Process.Signals = service.OnShutdownSignals(func(s os.Signal) {
		fmt.Println("[starship] received exit signal:", s)
		app.Stop() // NOTE: Stop has a call to CleanUp()
	})

	app.Process.WritePid(app.Config.Pid)
	app.ParseApplicationDirectories()
	//app.ParseUserDirectories()

	app.KV.NewCollection("users")
	return app
}

func (self *Application) CleanUp() error {
	fmt.Println("[starship] shutting down http server and closing session store")
	self.HTTPServer.Stop()
	fmt.Println("[starship] attempting to exit cleanly...")
	fmt.Println("[starship] closing the general key/value store")
	self.KV.Store.Close()
	fmt.Println("[starship] cleaning the pid file")
	self.Process.CleanPid()
	return nil
}
