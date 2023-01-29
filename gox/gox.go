package gox

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// gox 是一个简单的路由实现, 参考了 http.NewServeMux

var SupportMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost,
	http.MethodPut, http.MethodPatch, http.MethodDelete,
	http.MethodConnect, http.MethodOptions, http.MethodTrace}

type Gox struct {
	mutex  sync.RWMutex
	routes map[string]*route
}

func New() *Gox {
	return &Gox{
		mutex:  sync.RWMutex{},
		routes: make(map[string]*route),
	}
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

func (gx *Gox) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request), methods ...string) {
	gx.Handle(pattern, http.HandlerFunc(handler), methods...)
}

// Handle 添加路由/路由注册
// 规则：
/*
参考 gorilla/mux
"/products/{key}"
"/articles/{category}/"
"/articles/{category}/{id:[0-9]+}"
*/
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

	// 分割每一部分
	paths := strings.Split(pattern, "/")

	if gx.routes == nil {
		gx.routes = make(map[string]*route)
	}

	// 添加路由
	fmt.Println("添加路由", paths)
	r := route{h: handler, pattern: pattern, methods: methods, paths: paths}
	gx.routes[pattern] = &r
}

func (gx *Gox) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	for _, route := range gx.routes {
		fmt.Println("route.paths: ", route.paths)
		if ctx, ok := route.match(r.Context(), paths); ok {
			// http method 是否存
			exists := has(route.methods, r.Method)
			if !exists {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			route.h.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	}

	http.NotFound(w, r)
}

func has(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
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
