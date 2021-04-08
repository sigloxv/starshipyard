package framework

import (
	"fmt"
	"os"

	service "github.com/multiverse-os/starshipyard/framework/os/service"
	server "github.com/multiverse-os/starshipyard/framework/server"
)

// TODO: SysV, SystemD and Upstart init script creation and installation for
// intelligent and most importanlty secure defaults.
func (self *Application) StartAsDaemon() { service.Daemonize(func() { self.Start() }) }

func (self *Application) Start() {
	fmt.Println("[starship] starting the web application http server")
	self.Server[server.HTTP].Start()
	// TODO: Should hold open application until stop is called. Id like a better
	// way of holding the application open
	for {
	}
}

func (self *Application) Stop() {
	fmt.Println("[shipyard] initiating cleanup sequence, and stopping the starship process")
	// NOTE: Could not work with stores so they will automatically be handled below
	self.ShutdownFunctions()
	self.Process.CleanPid()
	self.Store.Close()
	os.Exit(0)
}

func (self *Application) Restart() {
	self.Server[server.HTTP].Stop()
	self.Server[server.HTTP].Start()
	fmt.Println("[shipyard] restarting the web application http server")
}
