package server

import (
	"fmt"
	"net/http"

	server "github.com/multiverse-os/starshipyard/framework/server"
)

type Route struct {
	Path string // So we can call them from views easily using *_path like root_path
	Name string
}

func LoadRoutes(server *server.Server) {
	server.DebutRoutes()
	server.UserRoutes()
}

func (self *Server) DebugRoutes() {
	//html.DefaultTemplate(self.FlashMessages, "starshipyard.io",
	self.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(self.Templates[DefaultTemplate].Render(r, "hello world").HTMLAsBytes())
	})
}

func (self *Server) UserRoutes() {
	// User signup, login, confirmation, and password recovery
	self.Router.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	self.Router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		sidCookie, err := r.Cookie("sid") // NOTE: cant get this far wtihout a sid, theoritically err=nil
		if err != nil {
			fmt.Println("this should never happen, should refresh ")
		} else {
			sid := sidCookie.Value
			uid := r.Form.Get("uid")
			password := r.Form.Get("password")
			// TODO: neds to populate flash messages
			self.UserLogin(sid, uid, password)
		}
		w.Write(self.Templates[DefaultTemplate].Render(r, "login failed").HTMLAsBytes())

	})
}
