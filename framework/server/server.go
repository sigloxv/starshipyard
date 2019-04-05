package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	config "github.com/multiverse-os/starshipyard/framework/config"
	kvstore "github.com/multiverse-os/starshipyard/framework/db/kvstore"
	template "github.com/multiverse-os/starshipyard/framework/html/template"
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
	Templates map[template.TemplateType]*template.Template
	Config    *config.Config
	Writer    http.ResponseWriter
	HTTP      *http.Server
	Router    router.Router
	Sessions  map[string]*Session
	sessions  *kvstore.KVStore
	cache     map[string]string
}

func New(config *config.Config) *Server {
	server := &Server{
		Router:   router.New(),
		Config:   config,
		Sessions: make(map[string]*Session),
		cache:    make(map[string]string),
	}
	server.LoadSessionStore()
	return server
}

// TODO: Split into logged in and not logged in sessions? Separate out sessions by roles?
// TODO: Migrate to using the higher level API for kvstore that is used by the
// framework.go file in the Init() func
func (self *Server) LoadSessionStore() {
	store, err := kvstore.New("db/sessions.db")
	if err != nil {
		panic(fmt.Sprintf("[fatal error] failed to open session DB: %v", err))
	}
	sessions, err := store.NewCollection("sessions")
	if err != nil {
		panic(fmt.Sprintf("[fatal error] failed to create 'sessions' kv collection: %v", err))
	}
	self.sessions = &kvstore.KVStore{
		Store: store,
		Collections: map[string]*kvstore.KeyValue{
			"sessions": sessions,
		},
	}
}

func (self *Server) Start() {
	//self.LoadMiddleware()
	//self.LoadTemplates()

	self.Router = starship.Router()

	self.HTTP = &http.Server{Addr: self.ListenAt(), Handler: self.Router}
	fmt.Println("[starship] http server listening on [ " + self.ListenAt() + " ]")
	self.HTTP.ListenAndServe()
}

func (self *Server) Stop() error {
	self.sessions.Store.Close()
	ctx, _ := context.WithTimeout(context.Background(), (15 * time.Second))
	if err := self.HTTP.Shutdown(ctx); err != nil {
		return fmt.Errorf("[error] failed to shutdown the http server:", err)
	}
	return nil
}

func (self *Server) ListenAt() string {
	return (self.Config.Address + ":" + strconv.Itoa(self.Config.Port))
}
