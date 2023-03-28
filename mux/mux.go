package mux

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

// muxMode 使用 mux package 的模式，预留着，暂没有复杂功能
var muxMode string = "debug" // debug ｜ dev ｜prod

// 支持的 HTTP Method
var listMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost,
	http.MethodPut, http.MethodPatch, http.MethodDelete,
	http.MethodConnect, http.MethodOptions, http.MethodTrace}

type mwHandler func(http.Handler) http.Handler

// ServeMux is an HTTP request multiplexer.
type ServeMux struct {
	routes        *Trie
	handlers      map[string]*router
	middlewares   []mwHandler
	notFound      http.Handler
	notAllowed    http.Handler
	methodOptions http.Handler
	mutex         sync.RWMutex
}

type router struct {
	pattern     string
	method      string
	handler     http.Handler
	middlewares []mwHandler
}

func NewServeMux() *ServeMux {
	srv := &ServeMux{
		routes:      NewTrie(),
		handlers:    make(map[string]*router),
		middlewares: make([]mwHandler, 0),
		mutex:       sync.RWMutex{},
	}

	srv.notFound = http.NotFoundHandler()

	srv.notAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	})

	srv.methodOptions = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	return srv
}

// Handle 注册路由
func (s *ServeMux) Handle(pattern string, handler http.Handler, methods ...string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}

	if handler == nil {
		panic("http: nil handler")
	}

	segments := strings.Split(pattern, "/")[1:]
	_ = s.routes.Insert(segments)
	if ok := s.routes.Insert(segments); ok {
		bytes, _ := json.Marshal(s.routes)
		log.Println("routes: ", string(bytes))
		panic(fmt.Sprintf("pattern is exists: %s", pattern))
	}
	if len(methods) == 0 {
		methods = listMethods
	}

	// 添加 handlers
	for _, m := range methods {
		r := &router{
			pattern:     pattern,
			method:      strings.ToUpper(m),
			handler:     handler,
			middlewares: make([]mwHandler, 0),
		}
		key := strings.ToUpper(m) + " " + pattern
		s.handlers[key] = r
	}
}

func (s *ServeMux) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request), methods ...string) {
	s.Handle(pattern, http.HandlerFunc(handler), methods...)
}

// ServeHTTP 实现 http.Handler 接口
func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")[1:]
	// 路由匹配
	ctx, match := s.routes.Match(r.Context(), segments)
	if !match {
		s.wrap(&router{handler: s.notFound}).ServeHTTP(w, r.WithContext(ctx))
		return
	}

	// 从 context 取出注册路由时的 pattern
	pattern, exists := ctx.Value(patternKey).(string)
	if !exists {
		s.wrap(&router{handler: s.notFound}).ServeHTTP(w, r.WithContext(ctx))
		return
	}

	// 查找 handler
	key := strings.ToUpper(r.Method) + " " + pattern
	logDebug("get handler key=%q", key)
	handler := s.handlers[key]
	if handler == nil {
		if ok := s.hasMedthod(pattern); ok {
			s.wrap(&router{handler: s.notAllowed}).ServeHTTP(w, r.WithContext(ctx))
			return
		}
		s.wrap(&router{handler: s.notFound}).ServeHTTP(w, r.WithContext(ctx))
		return
	}

	s.wrap(handler).ServeHTTP(w, r.WithContext(ctx))
}

func (s *ServeMux) SetNotFoundHandler(fn http.Handler) {
	s.notFound = fn
}

func (s *ServeMux) SetNoAllowedHandler(fn http.Handler) {
	s.notAllowed = fn
}

func (s *ServeMux) SetMethodOptionsHandler(fn http.Handler) {
	s.methodOptions = fn
}

// SetMuxMode 设置当前 ServeMux 模式，支持 debug ｜ dev ｜prod
func SetMuxMode(mode string) {
	if mode == "debug" || mode == "dev" || mode == "prod" {
		muxMode = mode
	}
	log.Printf("mode value must one of 'debug, dev, prod', but got %s", mode)
}

// wrap 包装，先执行中间件在执行逻辑处理
func (s *ServeMux) wrap(r *router) http.Handler {
	middlewares := append(s.middlewares, r.middlewares...)
	handler := r.handler
	for _, mw := range middlewares {
		handler = mw(handler)
	}
	return handler
}

// hasMedthod 判断pattern是否存在某个method
func (s *ServeMux) hasMedthod(pattern string) bool {
	for _, method := range listMethods {
		key := strings.ToUpper(method) + " " + pattern
		_, ok := s.handlers[key]
		return ok
	}
	return false
}

func logDebug(format string, v ...any) {
	if muxMode == "debug" {
		log.Printf(format+"\n", v...)
	}
}
