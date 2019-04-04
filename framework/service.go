package framework

import (
	"fmt"
	"os"

	service "github.com/multiverse-os/starshipyard/framework/service"
)

// TODO: Create a read-only database type we can use for things like the config DB
func (self *Application) StartAsDaemon() { service.Daemonize(func() { self.Start() }) }
func (self *Application) Start() {
	fmt.Println("[starship] starting the web application http server")
	self.HTTPServer.Start()
	// TODO: Should hold open application until stop is called. Id like a better
	// way of holding the application open
	for {
	}
}

func (self *Application) Stop() {
	fmt.Println("[shipyard] initiating the stop sequence")
	self.CleanUp()
	os.Exit(0)
}

func (self *Application) Restart() {
	self.HTTPServer.Stop()
	self.HTTPServer.Start()
	fmt.Println("[shipyard] restarting the web application http server")
}
