package service

import (
	"os"

	signal "github.com/multiverse-os/starshipyard/framework/os/service/signal"
)

func OnShutdownSignals(function func(os.Signal)) signal.Handler {
	return signal.ShutdownHandler(function)
}
