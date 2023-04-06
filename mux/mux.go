package mux

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
)

// TODO: 后面考虑是否要做的事情
// 支持 CORS
// 提取 Search Query

// muxMode 使用 mux package 的模式，预留着，暂没有复杂功能
var muxMode string = "debug" // debug ｜ dev ｜prod

type contextKey string

// 支持的 HTTP Method
var HTTPMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost,
	http.MethodPut, http.MethodPatch, http.MethodDelete,
	http.MethodConnect, http.MethodOptions, http.MethodTrace}

// MiddlewareFunc 中间件适配器函数签名
type MiddlewareFunc func(http.Handler) http.Handler

// ServeMux is an HTTP request multiplexer.
type ServeMux struct {
	NotFoundHandler  http.Handler
	MethodNotAllowed http.Handler
	MethodOptions    http.Handler
	isAllowed        bool // 是否支持 MethodNotAllowed，True开启，开启后，有一定的性能损耗
	routes           *Trie
	middlewares      []MiddlewareFunc
	mutex            sync.RWMutex
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

	mux.NotFoundHandler = http.NotFoundHandler()

	mux.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	})

	mux.MethodOptions = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	return mux
}

// Handle 注册路由总入口函数，所有的路由注册最终实现者
func (s *ServeMux) Handle(pattern string, handler http.Handler, methods ...string) Nodes {
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
	// logDebug("register route success", pattern)
	return trieNodes
}

// ServeHTTP 实现 http.Handler 接口
func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matchNode, exists := s.routes.Match(r)
	if !exists {
		for key := range s.routes.root.children {
			_, found := s.routes.Match(r, key)
			if found && s.isAllowed {
				s.wrap(s.MethodNotAllowed, s.middlewares).ServeHTTP(w, r)
				return
			}
		}

		s.wrap(s.NotFoundHandler, s.middlewares).ServeHTTP(w, r)
		return
	}

	// 设置 URL params 参数
	ctx := s.setParams(r.Context(), matchNode.fullpattern, r.URL.Path)

	r = r.WithContext(ctx)

	mws := []MiddlewareFunc{}
	mws = append(mws, s.middlewares...)
	mws = append(mws, matchNode.middlewares...)
	s.wrap(matchNode.handler, mws).ServeHTTP(w, r)
}

// Use 注册全局使用的中间件
func (s *ServeMux) Use(mws ...MiddlewareFunc) {
	s.middlewares = append(s.middlewares, mws...)
}

// OpenAllowed 开启支持 MethodNotAllowed
func (s *ServeMux) OpenAllowed() {
	s.isAllowed = true
}

// AllowedStatus 获取 isAllowed 状态
func (s *ServeMux) AllowedStatus() bool {
	return s.isAllowed
}

func (s *ServeMux) setParams(ctx context.Context, pattern string, url string) context.Context {
	// pattern: /api/:fish/fish/:id|^[a-zA-Z]+$/:age|^[0-9]+$
	// url      /api/feiyu/fish/99/GO/3
	if url == "/" {
		return ctx
	}

	ps := strings.Split(pattern, "/")[1:]
	us := strings.Split(url, "/")[1:]
	if len(ps) != len(us) {
		return ctx
	}
	for i, v := range ps {
		if strings.HasPrefix(v, ":") {
			paramKey, _, found := strings.Cut(v, "|")
			if found {
				paramKey = strings.Split(paramKey, ":")[1]
			} else {
				paramKey = strings.Split(v, ":")[1]
			}
			ctx = context.WithValue(ctx, contextKey(paramKey), us[i])
		}
	}

	return ctx
}

// wrap 包装，先执行中间件在执行逻辑处理
func (s *ServeMux) wrap(handler http.Handler, mws []MiddlewareFunc) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}

	return handler
}

// Param 获取param
func Param(r *http.Request, key string) string {
	val, ok := r.Context().Value(contextKey(key)).(string)
	// logDebug("Param(%q) is %t value=%q", key, ok, val)
	if !ok {
		return ""
	}
	return val
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

func (s *ServeMux) PrintTrieRoutes() {
	trieBytes, err := json.MarshalIndent(s.routes.root, "", " ")
	if err != nil {
		logDebug(err.Error())
	}
	logDebug(string(trieBytes))
}
