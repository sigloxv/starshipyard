package framework

// TODO: Need a handle on the functions so they can be removed
func (self *Application) AppendToShutdownProcess(exitFunction func() error) {
	self.Shutdown = append(self.Shutdown, exitFunction)
}

func (self *Application) ShutdownApplication() {
	for _, function := range self.Shutdown {
		function()
	}
}
