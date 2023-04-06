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
	router.GET("/", handlerEcho)
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

	http.ListenAndServe(":8888", router)
}

func handlerEcho(w http.ResponseWriter, r *http.Request) {
	params := map[string]string{}
	mux.Param(r, "name")
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
