package controllers

import (
	"fmt"
	"net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("root path")
	w.Write([]byte("hello world"))
}
