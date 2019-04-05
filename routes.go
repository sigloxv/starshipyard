package starship

import (
	"net/http"

	controllers "github.com/multiverse-os/starshipyard/controllers"
	router "github.com/multiverse-os/starshipyard/framework/server/router"
)

func Router() router.Router {
	r := router.New()

	r.Get("/", controllers.Root)

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	return r
}
