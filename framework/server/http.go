package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	router "github.com/multiverse-os/starshipyard/framework/server/router"

	starship "github.com/multiverse-os/starshipyard"
)

type HTTP struct {
	Config Config
	Writer http.ResponseWriter
	HTTP   *http.Server
	Router router.Router
	TLS    *TLS
	cache  map[string]string
}

func NewHTTP(address string, port int) Server {
	server := &HTTP{
		Router: router.New(),
		Config: Config{
			Address: address,
			Port:    port,
		},
	}
	return server
}

func (self *HTTP) IsRunning() bool {
	return true
}

func (self *HTTP) Start() error {
	self.Router = starship.Router()
	self.HTTP = &http.Server{Addr: self.ListeningAt(), Handler: self.Router}
	fmt.Println("[starship] http server listening on [ " + self.ListeningAt() + " ]")
	self.HTTP.ListenAndServe()
	return nil
}

func (self *HTTP) Stop() error {
	ctx, _ := context.WithTimeout(context.Background(), (15 * time.Second))
	if err := self.HTTP.Shutdown(ctx); err != nil {
		return fmt.Errorf("[error] failed to shutdown the http server:", err)
	}
	return nil
}

func (self *HTTP) ListeningAt() string {
	return (self.Config.Address + ":" + strconv.Itoa(self.Config.Port))
}
