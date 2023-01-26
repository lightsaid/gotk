package gox

import (
	"fmt"
	"net/http"
	"sync"
)

// gox 是一个简单的路由实现, 参考了 http.NewServeMux

type Gox struct {
	mutex  sync.RWMutex
	routes map[string]Route
}

type Route struct {
	h       http.Handler
	pattern string
}

func New() *Gox {
	return &Gox{
		mutex:  sync.RWMutex{},
		routes: make(map[string]Route),
	}
}

func (gx *Gox) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, ok := gx.routes[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	route.h.ServeHTTP(w, r)
}

func (gx *Gox) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	if handler == nil {
		panic("http: nil handler")
	}

	gx.Handle(pattern, http.HandlerFunc(handler))
}

func (gx *Gox) Handle(pattern string, handler http.Handler) {
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
		gx.routes = make(map[string]Route)
	}
	r := Route{h: handler, pattern: pattern}
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
