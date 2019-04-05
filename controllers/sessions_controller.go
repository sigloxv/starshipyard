package controllers

import (
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login controller")
	w.Write([]byte("GET login form goes here"))
}

func NewSession(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("POST login form to this controller and it then either forwards to login with errors or previous page logged in"))
}

// TODO: Maybe this should be in user controller, and separate user from session
// since session would include a lot of logic on its own
func Register(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET register form"))
}
