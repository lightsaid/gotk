package mux_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lightsaid/gotk/mux"
)

func helloMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(">>> Hello Middleware <<<")
		handler.ServeHTTP(w, r)
	})
}

func sinceMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		handler.ServeHTTP(w, r)
		log.Println(r.URL.Path, "Execute Time: ", time.Since(t))
	})
}

func TestXxx(t *testing.T) {
	s := mux.NewServeMux()

	// --- 注册路由 ---
	// 原来的
	s.GET("/v3/api/products/:name|^[a-zA-Z]+$", func(w http.ResponseWriter, r *http.Request) {
		log.Println("测试 RouteGroup 是否冲突")
	}).Use(sinceMiddleware)

	// 分组
	group := s.RouteGroup("/v3/api")
	group.Use(helloMiddleware)
	group.GET("/products/:id|^[0-9]+$", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.Method, r.URL.Path, "RouteGroup Testing....")
	}).Use(helloMiddleware, sinceMiddleware)

	// ---- 发请求 ----

	// 1. 原来的
	req, err := http.NewRequest("GET", "/v3/api/products/abc", nil)
	if err != nil {
		t.Error(err)
	}
	rsp := httptest.NewRecorder()
	s.ServeHTTP(rsp, req)
	fmt.Println(">>> ", rsp.Code)

	// 2. 分组
	req, err = http.NewRequest("GET", "/v3/api/products/999", nil)
	if err != nil {
		t.Error(err)
	}
	rsp = httptest.NewRecorder()
	s.ServeHTTP(rsp, req)
	fmt.Println(">>> ", rsp.Code)

}
