package framework

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	config "github.com/multiverse-os/starshipyard/framework/config"
	datastore "github.com/multiverse-os/starshipyard/framework/datastore"
	template "github.com/multiverse-os/starshipyard/framework/html/template"
	server "github.com/multiverse-os/starshipyard/framework/server"
	service "github.com/multiverse-os/starshipyard/framework/service"

	scramble "github.com/multiverse-os/scramble-key"
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
type Application struct {
	ScrambleKey scramble.Key
	Config      config.Config
	Process     *service.Process
	Directories ApplicationDirectories

	Shutdown []func() error

	Template map[template.TemplateType]*template.Template     // TODO: Should template data just be stored in a store?
	Store    map[datastore.DatastoreType]*datastore.Datastore // NOTE: Just store, but will make more sense when calling something from the map
	Server   map[server.ServerType]*server.Server
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

func Init(config config.Config) *Application {
	DropPriviledges()

	seedRandom()
	// NOTE: This is bare minimum validation and default fallbacks so that errors
	// are not thrown when setting up the application process signal handler, pid
	// control and other service functionality. An improved string/path validation
	// needs to be built ontop fo this basic functionality
	config = ValidateConfig(config)

	//wd, _ := os.Getwd()

	// TODr: SHould encapsualte all files into a embedded virtualFS so
	// transversals and similar attacks are within a virtual system that is
	// outside the actual FS or encapsulated so its segregated from the FS
	// preferably in reality stored in a BoltFS or similar type DB. Ideally in
	// blocks (crc32) or similar that can be scaled up by replicating across
	// harddisks to overcome IO bottlenecks
	app := &Application{
		ScrambleKey: scramble.GenerateKey(),
		Config:      config,
		Process:     service.ParseProcess(),
		Shutdown:    []func() error{},
		Template:    make(map[template.TemplateType]*template.Template),
		Store:       make(map[datastore.DatastoreType]*datastore.Datastore),
		Server:      make(map[server.ServerType]*server.Server),
	}
	app.Process.Signals = service.OnShutdownSignals(func(s os.Signal) {
		if s.String() == "interrupt" {
			fmt.Printf("\n")
		}
		fmt.Println("[starship] received exit signal:", s)
		app.Stop()
	})
	app.Process.WritePid(app.Config.Pid)

	//app.ParseApplicationrirectories()

	// TODO: Handle models
	// TODO: Load databases into Store map
	// TODO: Load HTTP server into server
	// mapaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa

	app.AppendToShutdownProcess(TestShutdownProcess)

	return app
}

func TestShutdownProcess() error {
	fmt.Println("SUCCESS! Shutdown process is running through appended functions!")
	return nil
}

func (self *Application) OpenKVStore(path string) (Close func()) {
	kv := datastore.OpenKVStore(path)
	return kv.Close
}
