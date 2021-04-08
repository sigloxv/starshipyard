package controllers

import (
	"fmt"
	"net/http"

	views "github.com/multiverse-os/starshipyard/views"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(views.Login().String()))
}

func NewSession(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("Entering controllers.Login()")

	uid := r.Form.Get("uid")
	fmt.Println("uid:", uid)

	pwd := r.Form.Get("pwd")
	fmt.Println("pwd:", pwd)

	fmt.Println("login controller")
	w.Write([]byte(views.Root().String()))
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(views.Root().String()))
}
