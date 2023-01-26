package gox

import (
	"fmt"
	"net/http"
	"sync"
)

// gox 是一个简单的路由实现, 参考了 http.NewServeMux

/*

package http

// Common HTTP methods.
//
// Unless otherwise noted, these are defined in RFC 7231 section 4.3.
const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" // RFC 5789
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

*/

var SupportMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost,
	http.MethodPut, http.MethodPatch, http.MethodDelete,
	http.MethodConnect, http.MethodOptions, http.MethodTrace}

type Gox struct {
	mutex  sync.RWMutex
	routes map[string]route
}

type route struct {
	h       http.Handler
	methods []string
	pattern string
}

func New() *Gox {
	return &Gox{
		mutex:  sync.RWMutex{},
		routes: make(map[string]route),
	}
}

func (gx *Gox) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, ok := gx.routes[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	if len(route.methods) == 0 {
		route.methods = SupportMethods
	}

	exists := has(route.methods, r.Method)
	if !exists {
		http.NotFound(w, r)
		return
	}

	route.h.ServeHTTP(w, r)
}

func has(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func (gx *Gox) GET(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	gx.Handle(pattern, http.HandlerFunc(handler), http.MethodGet)
}

func (gx *Gox) HEAD(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	gx.Handle(pattern, http.HandlerFunc(handler), http.MethodHead)
}

func (gx *Gox) POST(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	gx.Handle(pattern, http.HandlerFunc(handler), http.MethodPost)
}

func (gx *Gox) PUT(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	gx.Handle(pattern, http.HandlerFunc(handler), http.MethodPut)
}

func (gx *Gox) PATCH(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	gx.Handle(pattern, http.HandlerFunc(handler), http.MethodPatch)
}

func (gx *Gox) DELETE(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	gx.Handle(pattern, http.HandlerFunc(handler), http.MethodDelete)
}

func (gx *Gox) CONNECT(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	gx.Handle(pattern, http.HandlerFunc(handler), http.MethodConnect)
}

func (gx *Gox) OPTIONS(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	gx.Handle(pattern, http.HandlerFunc(handler), http.MethodOptions)
}

func (gx *Gox) TRACE(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	gx.Handle(pattern, http.HandlerFunc(handler), http.MethodTrace)
}

func (gx *Gox) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	// 交给 gx.Handle 验证
	// if handler == nil {
	// 	panic("http: nil handler")
	// }

	gx.Handle(pattern, http.HandlerFunc(handler))
}

// Handle 添加路由/路由注册
func (gx *Gox) Handle(pattern string, handler http.Handler, methods ...string) {
	gx.mutex.Lock()
	defer gx.mutex.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}

	if handler == nil {
		panic("http: nil handler")
	}

	if _, exist := gx.routes[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	if gx.routes == nil {
		gx.routes = make(map[string]route)
	}

	r := route{h: handler, pattern: pattern, methods: methods}
	gx.routes[pattern] = r
}

func Demo() {
	var handle = func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "demo")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)
	mux.Handle("/test", http.HandlerFunc(handle))
	// mux.ServeHTTP()
	// mux.Handler()
}
