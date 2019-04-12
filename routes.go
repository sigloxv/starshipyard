package starship

import (
	"net/http"

	controllers "github.com/multiverse-os/starshipyard/controllers"
	router "github.com/multiverse-os/starshipyard/framework/server/router"
)

// TODO: Want to generate these variables that provide helper links to routes
// automatically so we dont have to manually define them separately to the
// routes. These should be available in both the views and controllers, with a
// ""." import so they dont require a package name. So they can be used just
// like rails.
const (
	root_path       = "/"
	login_path      = "/login"
	public_key_path = "/public_key"
)

// TODO: To avoid the use of globals it is very likely we will need to pass the
// templates through here to have them be accessible in the controlers
func Router() router.Router {
	r := router.New()

	r.Get("/", controllers.Root)

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	return r
}

//import (
//	html "github.com/multiverse-os/starshipyard/framework/html"
//	template "github.com/multiverse-os/starshipyard/framework/html/template"
//)

// TODO: Building out rails-like path/route conviences preferably generated
// from a route defintion
// We want a system that:
//  * we have a function to init routes, providing us with these variables or
//  this but in a map.
//  * paths should be lowercase
// * outpput a title which chops of first / or last item of / split
