package mux

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

// muxMode 使用 mux package 的模式，预留着，暂没有复杂功能
var muxMode string = "debug" // debug ｜ dev ｜prod

// 支持的 HTTP Method
var HTTPMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost,
	http.MethodPut, http.MethodPatch, http.MethodDelete,
	http.MethodConnect, http.MethodOptions, http.MethodTrace}

// MiddlewareFunc 中间件适配器函数类型
type MiddlewareFunc func(http.Handler) http.Handler

// ServeMux is an HTTP request multiplexer.
type ServeMux struct {
	routes        *Trie
	middlewares   []MiddlewareFunc
	notFound      http.Handler
	notAllowed    http.Handler
	methodOptions http.Handler
	mutex         sync.RWMutex
	*routeGroup
}

// Nodes 暂存 Node 用来设置midddleware使用
type Nodes []*Node

// Use 给HTTP请求处理函数设置中间件
func (nodes Nodes) Use(mws ...MiddlewareFunc) {
	for _, node := range nodes {
		node.addMiddleware(mws...)
	}
}

// NewServeMux 创建一个多路复用HTTP路由器
func NewServeMux() *ServeMux {
	mux := &ServeMux{
		routes:      NewTrie(),
		middlewares: make([]MiddlewareFunc, 0),
		mutex:       sync.RWMutex{},
	}

	mux.routeGroup = &routeGroup{ServeMux: mux}

	mux.notFound = http.NotFoundHandler()

	mux.notAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	})

	mux.methodOptions = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	return mux
}

// Use 注册全局使用的中间件
func (s *ServeMux) Use(mws ...MiddlewareFunc) {
	s.middlewares = append(s.middlewares, mws...)
}

// Handle 注册路由总入口函数，所有的路由注册最终实现者
func (s *ServeMux) Handle(pattern string, handler http.Handler, methods ...string) Nodes {
	if s == nil {
		fmt.Println("s= ", s == nil)
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}

	if handler == nil {
		panic("http: nil handler")
	}

	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}

	pattern = cleanPath(pattern)

	// 如：将 /v1/api/tag/ 转化为 /v1/api/tag
	if pattern != "/" && strings.HasSuffix(pattern, "/") {
		pattern = pattern[0 : len(pattern)-1]
	}

	var trieNodes Nodes
	for _, method := range methods {
		newNode, err := s.routes.Insert(pattern, handler, strings.ToUpper(method))
		if err != nil {
			panic(err)
		}
		trieNodes = append(trieNodes, newNode)
	}
	// fmt.Println("register route success", pattern)
	// logDebug("register route success", pattern)
	return trieNodes
}

// ServeHTTP 实现 http.Handler 接口
func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matchNode, exists := s.routes.Match(r)
	if !exists {
		s.wrap(s.notFound, s.middlewares).ServeHTTP(w, r)
		return
	}

	mws := []MiddlewareFunc{}
	mws = append(mws, s.middlewares...)
	mws = append(mws, matchNode.middlewares...)
	s.wrap(matchNode.handler, mws).ServeHTTP(w, r)
}

// wrap 包装，先执行中间件在执行逻辑处理
func (s *ServeMux) wrap(handler http.Handler, mws []MiddlewareFunc) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}

	return handler
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

func cleanPath(path string) string {
	if !strings.Contains(path, "//") {
		return path
	}
	return strings.Replace(path, `//`, "/", -1)
}

func logDebug(format string, v ...any) {
	if muxMode == "debug" {
		log.Printf(format+"\n", v...)
	}
}
