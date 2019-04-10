package framework

import (
	"fmt"
)

// TODO: Need a handle on the functions so they can be removed
func (self *Application) AppendToShutdownProcess(exitFunction func() error) {
	self.Shutdown = append(self.Shutdown, exitFunction)
}

func (self *Application) ShutdownApplication() {
	fmt.Println("how many stored exit functions?", len(self.Shutdown))
	self.Shutdown[0]()

}
