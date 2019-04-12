package framework

// TODO: Need a handle on the functions so they can be removed
func (self *Application) AppendToShutdown(exitFunction func()) {
	self.Shutdown = append(self.Shutdown, exitFunction)
}

func (self *Application) ShutdownFunctions() {
	for _, function := range self.Shutdown {
		function()
	}
}
