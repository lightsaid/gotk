package gox

import (
	"net/http"
)

type routeGroup struct {
	gx          *Gox
	basePattern string
}

func (gx *Gox) RouteGroup(pattern string) *routeGroup {
	return &routeGroup{
		gx:          gx,
		basePattern: pattern,
	}
}

func (r *routeGroup) setPattern(pattern string) string {
	// 处理 /api/users//100
	path := r.basePattern + pattern
	size := len(r.basePattern)
	if size > 0 && len(pattern) > 0 {
		bPrefix := r.basePattern[:size-1]
		pPrefix := pattern[0:1]
		if bPrefix == "/" && pPrefix == "/" {
			path = r.basePattern + pattern[1:]
		}
	}
	// fmt.Println("group: ", path)
	return path
}

func (r *routeGroup) GET(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {

	return r.gx.Handle(r.setPattern(pattern), http.HandlerFunc(handler), http.MethodGet)
}

func (r *routeGroup) HEAD(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return r.gx.Handle(r.setPattern(pattern), http.HandlerFunc(handler), http.MethodHead)
}

func (r *routeGroup) POST(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return r.gx.Handle(r.setPattern(pattern), http.HandlerFunc(handler), http.MethodPost)
}

func (r *routeGroup) PUT(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return r.gx.Handle(r.setPattern(pattern), http.HandlerFunc(handler), http.MethodPut)
}

func (r *routeGroup) PATCH(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return r.gx.Handle(r.setPattern(pattern), http.HandlerFunc(handler), http.MethodPatch)
}

func (r *routeGroup) DELETE(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return r.gx.Handle(r.setPattern(pattern), http.HandlerFunc(handler), http.MethodDelete)
}

func (r *routeGroup) CONNECT(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return r.gx.Handle(r.setPattern(pattern), http.HandlerFunc(handler), http.MethodConnect)
}

func (r *routeGroup) OPTIONS(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return r.gx.Handle(r.setPattern(pattern), http.HandlerFunc(handler), http.MethodOptions)
}

func (r *routeGroup) TRACE(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return r.gx.Handle(r.setPattern(pattern), http.HandlerFunc(handler), http.MethodTrace)
}

func (r *routeGroup) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request), methods ...string) *routers {
	return r.gx.Handle(r.setPattern(pattern), http.HandlerFunc(handler), methods...)
}
