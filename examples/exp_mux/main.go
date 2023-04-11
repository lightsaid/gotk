package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/lightsaid/gotk/mux"
)

func main() {
	router := mux.NewServeMux()
	router.GET("/", handlerEcho).Use(helloMiddleware)
	router.GET("/api/", handlerEcho)
	router.GET("/api/:name", handlerEcho)
	router.GET("/api/:name/:id|^[0-9]+", handlerEcho)
	router.GET("/api/:cat/cat", handlerEcho)
	router.GET("/api/:dog/dog", handlerEcho)
	router.GET("/api/:fish/fish", handlerEcho)
	router.GET("/api/:fish/fish/:id|^[a-zA-Z]+$", handlerEcho)
	router.GET("/api/:fish/fish/:id|^[a-zA-Z]+$/:age|^[0-9]+$", handlerEcho)

	router.Use(sinceMiddleware, helloMiddleware)
	router.OpenAllowed()

	// 路由分组
	group := router.RouteGroup("/v1/auth")

	// 局部中间件，仅对这一组路由起效
	group.Use(sinceMiddleware)

	group.POST("/login", handlerEcho).Use(helloMiddleware)

	// 支持多 method, 如果不指定 Method, 默认支持所有
	group.HandleFunc("/profile", handlerEcho, http.MethodGet, http.MethodPost)

	// 静态资源
	// fs := http.FileServer(http.Dir("./static"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))

	router.Static("/static/", "./static")

	// router.PrintTrieRoutes()

	http.ListenAndServe(":8888", router)
}

func handlerEcho(w http.ResponseWriter, r *http.Request) {
	params := map[string]string{}

	params["name"] = mux.Param(r, "name")
	params["id"] = mux.Param(r, "id")
	params["cat"] = mux.Param(r, "cat")
	params["dog"] = mux.Param(r, "dog")
	params["fish"] = mux.Param(r, "fish")
	params["age"] = mux.Param(r, "age")

	json.NewEncoder(w).Encode(params)
}

func helloMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hello World")
		handler.ServeHTTP(w, r)
	})
}

func sinceMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf("%s %s %s \n", r.Method, r.URL.Path, time.Since(t))
	})
}
