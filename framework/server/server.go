package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	config "github.com/multiverse-os/starshipyard/framework/config"
	router "github.com/multiverse-os/starshipyard/framework/server/router"

	starship "github.com/multiverse-os/starshipyard"
)

// NOTE: Sessions are low level, especially in the Starship Yard model where we
// automatically assign a session to every visitor to provide DOS protection and
// other low-level connection related functionality.
// TODO: Need a way to add routes on-the-fly
// TODO: Need to add attribute for storing middleware, routes, and templates so
// it can be loaded in this file but declared at the higher levels to segregate
// the framework code from the application code.
type Server struct {
	Config *config.Config
	Writer http.ResponseWriter
	HTTP   *http.Server
	Router router.Router
	cache  map[string]string
}

func New(config *config.Config) *Server {
	server := &Server{
		Router: router.New(),
		Config: config,
		cache:  make(map[string]string),
	}
	return server
}

// TODO: Split into logged in and not logged in sessions? Separate out sessions by roles?
// TODO: Migrate to using the higher level API for kvstore that is used by the
// framework.go file in the Init() func
func (self *Server) Start() {

	self.Router = starship.Router()

	self.HTTP = &http.Server{Addr: self.ListenAt(), Handler: self.Router}
	fmt.Println("[starship] http server listening on [ " + self.ListenAt() + " ]")
	self.HTTP.ListenAndServe()
}

func (self *Server) Stop() error {
	ctx, _ := context.WithTimeout(context.Background(), (15 * time.Second))
	if err := self.HTTP.Shutdown(ctx); err != nil {
		return fmt.Errorf("[error] failed to shutdown the http server:", err)
	}
	return nil
}

func (self *Server) ListenAt() string {
	return (self.Config.Address + ":" + strconv.Itoa(self.Config.Port))
}
