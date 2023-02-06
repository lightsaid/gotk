package gox

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

// gox 是一个简单的路由实现, 参考 http.NewServeMux、 gorilla/mux、gin、flow...

var SupportMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost,
	http.MethodPut, http.MethodPatch, http.MethodDelete,
	http.MethodConnect, http.MethodOptions, http.MethodTrace}

type Gox struct {
	NotFound         http.Handler
	MethodNotAllowed http.Handler
	MethodOptions    http.Handler
	mutex            sync.RWMutex
	routes           RouteTrie
	router           map[string]*router

	// 存储中间件, 采用具名中间方式
	middlewares map[string]func(http.Handler) http.Handler
	// 全局使用的中间件的 key
	mwKeys []string
}

type router struct {
	handler http.Handler
	methods []string
	pattern string
	// 执行请求之前的中间件key
	beforeMWKeys []string
	// // 执行请求之后的中间件key
	afterMWKeys []string
}

func New() *Gox {
	return &Gox{
		NotFound: http.NotFoundHandler(),
		MethodNotAllowed: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}),
		MethodOptions: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
		mutex:       sync.RWMutex{},
		routes:      NewRouteTrie(),
		router:      make(map[string]*router),
		middlewares: make(map[string]func(http.Handler) http.Handler),
	}
}

// Register 注册中间件
func (gx *Gox) Register(key string, fn func(http.Handler) http.Handler) {
	gx.middlewares[key] = fn
}

// Global 全局使用的中间件，按照 keys 顺序执行
func (gx *Gox) Global(keys ...string) {
	gx.mwKeys = keys
}

type routers struct {
	routes []*router
}

// Before 一个路由执行之前使用的中间件，按照 keys 顺序执行
func (rs *routers) Before(keys ...string) *routers {
	for index := range rs.routes {
		rs.routes[index].beforeMWKeys = keys
	}

	return rs
}

// After 一个路由执行之后使用的中间件，按照 keys 顺序执行
func (rs *routers) After(keys ...string) *routers {
	for index := range rs.routes {
		rs.routes[index].afterMWKeys = keys
	}

	return rs
}

func (gx *Gox) Handle(pattern string, handler http.Handler, methods ...string) *routers {
	gx.mutex.Lock()
	defer gx.mutex.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}

	if handler == nil {
		panic("http: nil handler")
	}

	// 同一个线路可以注册为GET、POST、PUT...
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
	var rs routers
	for _, method := range methods {
		r := &router{handler: handler, methods: methods, pattern: pattern}
		rs.routes = append(rs.routes, r)
		gx.router[pattern+"#"+method] = r
	}

	return &rs
}

func (gx *Gox) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	parts := strings.Split(req.URL.Path, "/")[1:]
	ctx, ok := gx.routes.Match(req.Context(), parts)
	if !ok {
		gx.wrap(gx.NotFound, gx.mwKeys).ServeHTTP(w, req)
		return
	}
	newReq := req.WithContext(ctx)
	key, exists := ctx.Value(trieNodeKey).(string)
	if !exists {
		gx.wrap(gx.NotFound, gx.mwKeys).ServeHTTP(w, newReq)
		return
	}
	r, exists := gx.router[key+"#"+req.Method]
	if !exists {
		gx.wrap(gx.NotFound, gx.mwKeys).ServeHTTP(w, newReq)
		return
	}
	if r.pattern == key {
		if has(r.methods, req.Method) {
			// 执行全局中间件和before中间件
			var keys []string
			keys = append(keys, gx.mwKeys...)
			keys = append(keys, r.beforeMWKeys...)
			gx.wrap(r.handler, keys).ServeHTTP(w, newReq)

			// after 中间件
			gx.wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), r.afterMWKeys).ServeHTTP(w, newReq)
			return
		} else {
			if newReq.Method == http.MethodOptions {
				gx.wrap(gx.MethodOptions, gx.mwKeys).ServeHTTP(w, newReq)
				return
			}
			gx.wrap(gx.MethodNotAllowed, gx.mwKeys).ServeHTTP(w, newReq)
		}
	}

	gx.wrap(gx.NotFound, gx.mwKeys).ServeHTTP(w, newReq)
}

// wrap 执行中间件
func (gx *Gox) wrap(handler http.Handler, keys []string) http.Handler {
	for i := len(keys) - 1; i >= 0; i-- {
		fmt.Printf("len(keys): %d, i: %d \n", len(keys), i)
		// handler = gx.middlewares[keys[i]](handler)
		if fn, ok := gx.middlewares[keys[i]]; ok {
			handler = fn(handler)
		} else {
			log.Printf("middleware key: %s not exists\n", keys[i])
		}

	}
	return handler
}

func has(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
