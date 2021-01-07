package router

import (
	"net/http"
)

type ChainHandler struct {
	Middlewares Middlewares
	Endpoint    http.Handler
	// Private
	chain http.Handler
}

func Chain(middlewares ...func(http.Handler) http.Handler) Middlewares {
	return Middlewares(middlewares)
}

func (self Middlewares) Handler(h http.Handler) http.Handler {
	return &ChainHandler{
		Middlewares: self,
		Endpoint:    h,
		chain:       chain(self, h),
	}
}

func (self Middlewares) HandlerFunc(h http.HandlerFunc) http.Handler {
	return &ChainHandler{
		Middlewares: self,
		Endpoint:    h,
		chain:       chain(self, h),
	}
}

func (self *ChainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.chain.ServeHTTP(w, r)
}

func chain(middlewares []func(http.Handler) http.Handler, endpoint http.Handler) http.Handler {
	if len(middlewares) == 0 {
		return endpoint
	}

	// Wrap the end handler with the middleware chain
	h := middlewares[len(middlewares)-1](endpoint)
	for i := len(middlewares) - 2; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}
