package router

import (
	"context"
	"net"
	"net/http"
	"strings"
)

var (
	RouteCtxKey = &contextKey{"RouteContext"}
)

type Context struct {
	// Public
	Routes        Routes
	RoutePath     string
	RouteMethod   string
	RoutePatterns []string
	URLParams     RouteParams
	// Private
	routePattern     string
	routeParams      RouteParams
	methodNotAllowed bool
}

func NewRouteContext() *Context { return &Context{} }

func (self *Context) Reset() {
	self.Routes = nil
	self.RoutePath, self.RouteMethod, self.routePattern = ""
	self.RoutePatterns = self.RoutePatterns[:0]
	self.URLParams.Keys = self.URLParams.Keys[:0]
	self.URLParams.Values = self.URLParams.Values[:0]

	self.routeParams.Keys = self.routeParams.Keys[:0]
	self.routeParams.Values = self.routeParams.Values[:0]
	self.methodNotAllowed = false
}

func (self *Context) URLParam(key string) string {
	for k := len(self.URLParams.Keys) - 1; k >= 0; k-- {
		if self.URLParams.Keys[k] == key {
			return self.URLParams.Values[k]
		}
	}
	return ""
}

func (self *Context) RoutePattern() string {
	routePattern := strings.Join(self.RoutePatterns, "")
	return strings.Replace(routePattern, "/*/", "/", -1)
}

func RouteContext(ctx context.Context) *Context {
	return ctx.Value(RouteCtxKey).(*Context)
}

func URLParam(r *http.Request, key string) string {
	if rctx := RouteContext(r.Context()); rctx != nil {
		return rctx.URLParam(key)
	}
	return ""
}

func URLParamFromCtx(ctx context.Context, key string) string {
	if rctx := RouteContext(ctx); rctx != nil {
		return rctx.URLParam(key)
	}
	return ""
}

type RouteParams struct {
	Keys   []string
	Values []string
}

// TODO: We have a key/value store here, this may be a place of optimization
func (self *RouteParams) Add(key, value string) {
	(*self).Keys = append((*self).Keys, key)
	(*self).Values = append((*self).Values, value)
}

func ServerBaseContext(baseCtx context.Context, h http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		baseCtx := baseCtx
		if v, ok := ctx.Value(http.ServerContextKey).(*http.Server); ok {
			baseCtx = context.WithValue(baseCtx, http.ServerContextKey, v)
		}
		if v, ok := ctx.Value(http.LocalAddrContextKey).(net.Addr); ok {
			baseCtx = context.WithValue(baseCtx, http.LocalAddrContextKey, v)
		}

		h.ServeHTTP(w, r.WithContext(baseCtx))
	})
	return fn
}

type contextKey struct {
	name string
}

func (self *contextKey) String() string {
	return "router context value " + self.name
}
