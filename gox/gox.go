package gox

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// gox 是一个简单的路由实现, 参考 http.NewServeMux、 gorilla/mux、gin、flow...

var SupportMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost,
	http.MethodPut, http.MethodPatch, http.MethodDelete,
	http.MethodConnect, http.MethodOptions, http.MethodTrace}

type Gox struct {
	mutex  sync.RWMutex
	routes RouteTrie
	router map[string]*router
}

type router struct {
	handler http.Handler
	methods []string
	pattern string
}

func New() *Gox {
	return &Gox{
		mutex:  sync.RWMutex{},
		routes: NewRouteTrie(),
		router: make(map[string]*router),
	}
}

func Param(r *http.Request, key string) string {
	val, ok := r.Context().Value(contextKey(key)).(string)
	if !ok {
		return ""
	}
	return val
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

func (gx *Gox) Handle(pattern string, handler http.Handler, methods ...string) {
	gx.mutex.Lock()
	defer gx.mutex.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}

	if handler == nil {
		panic("http: nil handler")
	}

	// 同一个路由可以注册为GET、POST、PUT...
	// if _, exist := gx.router[pattern]; exist {
	// 	panic("http: multiple registrations for " + pattern)
	// }

	// 分割每一部分
	parts := strings.Split(pattern, "/")[1:]
	// if ok := gx.routes.Register(parts); !ok {
	// 	panic("http: multiple registrations for " + pattern)
	// }
	// 允许重复注册
	gx.routes.Register(parts)

	// 添加路由
	// fmt.Println("添加路由", parts)
	for _, method := range methods {
		gx.router[pattern+"#"+method] = &router{handler: handler, methods: methods, pattern: pattern}
	}
}

func (gx *Gox) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	parts := strings.Split(req.URL.Path, "/")[1:]
	ctx, ok := gx.routes.Match(req.Context(), parts)
	if !ok {
		http.NotFound(w, req)
		return
	}
	key, exists := ctx.Value(trieNodeKey).(string)
	if !exists {
		http.NotFound(w, req)
		return
	}
	r, exists := gx.router[key+"#"+req.Method]
	if !exists {
		http.NotFound(w, req)
		return
	}
	if r.pattern == key && has(r.methods, req.Method) {
		r.handler.ServeHTTP(w, req.WithContext(ctx))
		fmt.Println(">>> time: ", time.Since(t), r.pattern, req.URL.Path)
		return
	} else {
		http.NotFound(w, req)
		return
	}
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
