package controllers

import (
	"fmt"
	"net/http"

	//template "github.com/multiverse-os/starshipyard/framework/html/template"
	views "github.com/multiverse-os/starshipyard/views"
)

func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("root path")
	w.Write([]byte(views.Root().String()))
}
