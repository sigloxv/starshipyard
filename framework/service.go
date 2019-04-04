package framework

import (
	"fmt"
	"os"
	//service "github.com/multiverse-os/starshipyard/framework/service"
)

// TODO: Create a read-only database type we can use for things like the config DB
//func (self *Application) StartAsDaemon() { service.Daemonize(func() { self.Start() }) }
func (self *Application) Start() {
	fmt.Println("[starship] starting the web application http server")

	fmt.Println("app.Process.PidFile in start:", self.Process)

	fmt.Println("app.Process.Path in start:", self.Process.PidFile.Path)
	fmt.Println("app.Process.Pid in start:", self.Process.PidFile.Pid)
	self.HTTPServer.Start()
	// TODO: Should hold open application until stop is called. Id like a better
	// way of holding the application open
	for {
	}
	fmt.Println("its not holding things open :(")
}

func (self *Application) Stop() {
	fmt.Println("[shipyard] stopping the web application http server")
	self.HTTPServer.Stop()
	self.CleanUp()
	os.Exit(0)
}

func (self *Application) Restart() {
	self.HTTPServer.Stop()
	self.HTTPServer.Start()
	fmt.Println("[shipyard] restarting the web application http server")
}
