package controllers

import (
	"fmt"
	"net/http"

	//template "github.com/multiverse-os/starshipyard/framework/html/template"
	views "github.com/multiverse-os/starshipyard/views"
)

func Template(yield string) string {
	return ""
}

// TODO: In rails, we define the template in the controller. To acheive similar
// functionality, we will have functions that can be called from the controller
// to load and render desired template
func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("root path")

	// TODO: In the near future change this to render to and specify the format to
	// be more similar to rails for rails developers
	w.Write([]byte(views.Root().String()))
}
