package framework

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	config "github.com/multiverse-os/starshipyard/framework/config"
	datastore "github.com/multiverse-os/starshipyard/framework/datastore"
	filesystem "github.com/multiverse-os/starshipyard/framework/os/filesystem"
	service "github.com/multiverse-os/starshipyard/framework/os/service"
	"github.com/multiverse-os/starshipyard/framework/os/service/signal"
	server "github.com/multiverse-os/starshipyard/framework/server"

	scramble "github.com/multiverse-os/scramble-key"
)

// Aliasing for simplicity, enables the API to be app.Server[HTTP] instead of
// app.Server[server.HTTP] without calling the entire server package with `.`
var (
	HTTP = server.HTTP
	//KV   = datastore.KV
	//ObjectStore = datastore.ObjectStore
	//Cache       = datastore.Cache
	//Session     = datastore.Session
)

// NOTE: Concept: we want to be able to run multiple applications in a given
// instance. This would likely be defined by a ruby-like script config that
// defines what domains go where, reverse and inverting proxy settings, etc
// TODO: Starship Yard is meant to function as both the web application and
// include the firewall/reverse-proxy. This means we need to be able to register
// multiple end-points.
// TODO: Migrate to less pointer use, we should only be using pointers when
// absolutely necessary to do so, or when there is an obvious benefit. It should
// absolutely not be the default, just like globals.
// TODO: Need to persist and load scramble key
// TODO: Build a string function to provide a nice ouput with all necesary
// information

//var Store datastore.KV // NOTE: Just store, but will make more sense when calling something from the map

type Application struct {
	ScrambleKey scramble.Key
	Config      config.Settings
	Process     *service.Process
	Directories filesystem.ApplicationDirectories
	Shutdown    []func()
	Store       *datastore.KV // NOTE: Just store, but will make more sense when calling something from the map
	Server      map[server.ServerType]server.Server
}

func seedRandom() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func DropPriviledges() {
	if service.IsRootUser() {
		fmt.Println("[starship] running internet facing servers as root is very dangerous, run as an unpriviledged user")
		os.Exit(1)
		// TODO: Lets drop priviledges
	}
}

func Init(config config.Settings) *Application {
	DropPriviledges()
	seedRandom()
	// NOTE: This is bare minimum validation and default fallbacks so that errors
	// are not thrown when setting up the application process signal handler, pid
	// control and other service functionality. An improved string/path validation
	// needs to be built ontop fo this basic functionality
	config = ValidateConfig(config)

	// TODO: SHould encapsualte all files into a embedded virtualFS so
	// transversals and similar attacks are within a virtual system that is
	// outside the actual FS or encapsulated so its segregated from the FS
	// preferably in reality stored in a BoltFS or similar type DB. Ideally in
	// blocks (crc32) or similar that can be scaled up by replicating across
	// harddisks to overcome IO bottlenecks
	kvStore, err := datastore.OpenKV("db/kvstore")
	if err != nil {
		panic(errors.New("[error] failed to open leveldb datastore:" + err.Error()))
	}

	// TODO: Load Session Store
	// TODO: Load Scramble Key from config folder or generate a new one if there
	// is none. This will be used in our session encryption, and various other
	// aspects of the design.

	app := &Application{
		ScrambleKey: scramble.GenerateKey(),
		Config:      config,
		Process:     service.ParseProcess(),
		Store:       kvStore,
		Server:      make(map[server.ServerType]server.Server),
	}
	// Process Information Parsing and Long running Linux service initialization
	/////////////////////////////////////////////////////////////////////r///////
	app.Process.Signals = service.OnShutdownSignals(func(shutdownSignal os.Signal) {
		// TODO: No string comparisons, or at least first start with length compare
		switch shutdownSignal {
		case signal.Interrupt:
			fmt.Printf("\n")
		}
		fmt.Println("[starship] received exit signal:", shutdownSignal)
		app.Stop()
	})
	app.Process.WritePid(app.Config.Pid)

	//app.ParseApplicationrirectories()
	app.AppendToShutdown(TestShutdownProcess)

	//app.AppendToShutdownProcess(kv.Close)
	// TODO:  Initialize a cache DB with TTL or something similar

	// TODO: Support unix socket connections. Then several servers listening on
	// unix sockets can be proxied with a single server listening on the port and
	// address

	app.Server[HTTP] = server.NewHTTP(config.Address, config.Port)

	return app
}

func TestShutdownProcess() {
	fmt.Println("SUCCESS! Shutdown process is running through appended functions!")
}
