package framework

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	config "github.com/multiverse-os/starshipyard/framework/config"
	database "github.com/multiverse-os/starshipyard/framework/database"
	template "github.com/multiverse-os/starshipyard/framework/html/template"
	server "github.com/multiverse-os/starshipyard/framework/server"
	service "github.com/multiverse-os/starshipyard/framework/service"

	scramble "github.com/multiverse-os/scramble-key"
)

// NOTE: Concept: we want to be able to run multiple applications in a given
// instance. This would likely be defined by a ruby-like script config that
// defines what domains go where, reverse and inverting proxy settings, etc
//type Domain struct {
//	Name        string
//	Subdomains  []string
//	Certificate string
//}

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
	Config      *config.Config
	Process     *service.Process
	HTTPServer  *server.Server
	Templates   map[template.TemplateType]*template.Template
	ScrambleKey scramble.Key
	KV          *database.KV
	Sessions    *database.KV
	Directories ApplicationDirectories
}

func seedRandom() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func Init(config *config.Config) *Application {
	seedRandom()

	if service.IsRootUser() {
		fmt.Println("[starship] running internet facing servers as root is very dangerous, run as an unpriviledged user")
		os.Exit(1)
		// TODO: Lets drop priviledges
	}

	//wd, _ := os.Getwd()

	// NOTE: This is bare minimum validation and default fallbacks so that errors
	// are not thrown when setting up the application process signal handler, pid
	// control and other service functionality. An improved string/path validation
	// needs to be built ontop fo this basic functionality
	if len(config.Pid) == 0 {
		config.Pid = "tmp/pids/starship.pid"
		fmt.Println("config.Pid set to because it was blank:", config.Pid)
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

	// TODO: Can put the *.db in a memFS for a more transient pure memory DB
	// TODO: SHould encapsualte all files into a embedded virtualFS so
	// transversals and similar attacks are within a virtual system that is
	// outside the actual FS or encapsulated so its segregated from the FS
	// preferably in reality stored in a BoltFS or similar type DB. Ideally in
	// blocks (crc32) or similar that can be scaled up by replicating across
	// harddisks to overcome IO bottlenecks
	app := &Application{
		Config:      config,
		HTTPServer:  server.New(config),
		KV:          database.InitKV("db/kv.db"),
		Sessions:    database.InitKV("db/sessions.db").WithCollection("sessions"),
		ScrambleKey: scramble.GenerateKey(),
		Process:     service.ParseProcess(),
		Templates:   make(map[template.TemplateType]*template.Template),
	}

	app.Process.Signals = service.OnShutdownSignals(func(s os.Signal) {
		if s.String() == "interrupt" {
			fmt.Printf("\n")
		}
		fmt.Println("[starship] received exit signal:", s)
		app.Stop()
	})

	fmt.Println("[starship] writing pid:", app.Config.Pid)
	app.Process.WritePid(app.Config.Pid)
	app.ParseApplicationDirectories()
	//app.ParseUserDirectories()

	// TODO: THESE NEED TO BE NOT DEFINED HERE, they MUST be defined by the models
	// in the models folder, not in the framework. Infact  ALL application
	// settings need to be migrated out of this and only overridable defaults or
	// security oriented decisions should be in the framework portion of the
	// codebase
	app.KV.NewCollection("users")

	// TODO: Load controllers, models, etc

	return app
}
