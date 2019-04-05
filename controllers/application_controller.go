package controllers

import (
	"fmt"
	"net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("root path")
	w.Write([]byte("hello world"))
}

//func Routes() map[string]route.Route {
//	return map[string]route.Route{
//		"root": route.Route{
//			Pattern: "/",
//			Method:  route.Get,
//			Controller: func(w http.ResponseWriter, r *http.Request) {
//				fmt.Println("root path!")
//				//w.Write(self.Templates[DefaultTemplate].Render(r, "hello world").HTMLAsBytes())
//				w.Write([]byte("hello world"))
//			},
//		},
//		"new_user_session": route.Route{
//			Pattern: "/users/login",
//			Method:  route.Get,
//			Controller: func(w http.ResponseWriter, r *http.Request) {
//				w.Write([]byte("new user session"))
//			},
//		},
//		"user_session": route.Route{
//			Pattern: "/users/login",
//			Method:  route.Post,
//			Controller: func(w http.ResponseWriter, r *http.Request) {
//				//r.ParseForm()
//				//sidCookie, err := r.Cookie("sid") // NOTE: cant get this far wtihout a sid, theoritically err=nil
//				//if err != nil {
//				//	fmt.Println("this should never happen, should refresh ")
//				//} else {
//				//	//sid := sidCookie.Value
//				//	//uid := r.Form.Get("uid")
//				//	//password := r.Form.Get("password")
//				//	// TODO: neds to populate flash messages
//				//	//self.UserLogin(sid, uid, password)
//				//}
//				//w.Write(self.Templates[DefaultTemplate].Render(r, "login failed").HTMLAsBytes())
//				w.Write([]byte("log user in"))
//			},
//		},
//	}
//}
