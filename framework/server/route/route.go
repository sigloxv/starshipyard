package route

import (
	"net/http"
	"strings"
)

type Route struct {
	Pattern    string
	Method     Method
	Controller http.HandlerFunc
}

type Method int

const (
	Connect Method = iota
	Delete
	Get
	Head
	Options
	Patch
	Post
	Put
	Trace
)

func MarshallMethod(method string) Method {
	switch strings.ToLower(method) {
	case Connect.String():
		return Connect
	case Delete.String():
		return Delete
	case Head.String():
		return Head
	case Options.String():
		return Options
	case Patch.String():
		return Patch
	case Post.String():
		return Post
	case Put.String():
		return Put
	case Trace.String():
		return Trace
	default: // case Get.String():
		return Get
	}
}

func (self Method) String() string {
	switch self {
	case Connect:
		return "connect"
	case Delete:
		return "delete"
	case Get:
		return "get"
	case Head:
		return "head"
	case Options:
		return "options"
	case Patch:
		return "patch"
	case Post:
		return "post"
	case Put:
		return "put"
	case Trace:
		return " Trace"
	default:
		return "invalid"
	}
}
