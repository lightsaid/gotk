package mux

import (
	"net/http"
)

type routeGroup struct {
	prefix      string
	middlewares []MiddlewareFunc
	*ServeMux
}

// RouteGroup 路由分组
func (s *ServeMux) RouteGroup(pattern string) routeGroup {
	group := routeGroup{prefix: pattern, ServeMux: s}

	return group
}

func (r *routeGroup) Use(mws ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, mws...)
}

func (r *routeGroup) handle(pattern string, handler func(w http.ResponseWriter, r *http.Request), methods ...string) Nodes {
	nodes := r.Handle(pattern, http.HandlerFunc(handler), methods...)
	nodes.Use(r.middlewares...)
	return nodes
}

func (r *routeGroup) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request), methods ...string) Nodes {
	return r.handle(r.prefix+pattern, http.HandlerFunc(handler), methods...)
}

func (r *routeGroup) GET(pattern string, handler func(w http.ResponseWriter, r *http.Request)) Nodes {
	return r.handle(r.prefix+pattern, http.HandlerFunc(handler), http.MethodGet)
}

func (r *routeGroup) POST(pattern string, handler func(w http.ResponseWriter, r *http.Request)) Nodes {
	return r.handle(r.prefix+pattern, http.HandlerFunc(handler), http.MethodPost)
}

func (r *routeGroup) PUT(pattern string, handler func(w http.ResponseWriter, r *http.Request)) Nodes {
	return r.handle(r.prefix+pattern, http.HandlerFunc(handler), http.MethodPut)
}

func (r *routeGroup) PATCH(pattern string, handler func(w http.ResponseWriter, r *http.Request)) Nodes {
	return r.handle(r.prefix+pattern, http.HandlerFunc(handler), http.MethodPatch)
}

func (r *routeGroup) HEAD(pattern string, handler func(w http.ResponseWriter, r *http.Request)) Nodes {
	return r.handle(r.prefix+pattern, http.HandlerFunc(handler), http.MethodHead)
}

func (r *routeGroup) DELETE(pattern string, handler func(w http.ResponseWriter, r *http.Request)) Nodes {
	return r.handle(r.prefix+pattern, http.HandlerFunc(handler), http.MethodDelete)
}

func (r *routeGroup) CONNECT(pattern string, handler func(w http.ResponseWriter, r *http.Request)) Nodes {
	return r.handle(r.prefix+pattern, http.HandlerFunc(handler), http.MethodConnect)
}

func (r *routeGroup) OPTIONS(pattern string, handler func(w http.ResponseWriter, r *http.Request)) Nodes {
	return r.handle(r.prefix+pattern, http.HandlerFunc(handler), http.MethodOptions)
}

func (r *routeGroup) TRACE(pattern string, handler func(w http.ResponseWriter, r *http.Request)) Nodes {
	return r.handle(r.prefix+pattern, http.HandlerFunc(handler), http.MethodTrace)
}

// PPT 以 Post 和 Put 方法注册路由
func (r *routeGroup) PPT(pattern string, handler func(w http.ResponseWriter, r *http.Request)) Nodes {
	return r.handle(r.prefix+pattern, http.HandlerFunc(handler), http.MethodPost, http.MethodPut)
}

// PGT 以 Post 和 Get 方法注册路由
func (r *routeGroup) PGT(pattern string, handler func(w http.ResponseWriter, r *http.Request)) Nodes {
	return r.handle(r.prefix+pattern, http.HandlerFunc(handler), http.MethodPost, http.MethodGet)
}
