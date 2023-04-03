package mux

import (
	"log"
	"net/http"
	"time"
)

func loggerMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("INFO %s %s\n", r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
	})
}

func sinceMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		handler.ServeHTTP(w, r)
		log.Println(r.URL.Path, " exec time: ", time.Since(t))
	})
}

func helloMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(">>> Hello Middleware <<<")
		handler.ServeHTTP(w, r)
	})
}

func HelloMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(">>> Hello Middleware <<<")
		handler.ServeHTTP(w, r)
	})
}

func SinceMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		handler.ServeHTTP(w, r)
		log.Println(r.URL.Path, "Execute Time: ", time.Since(t))
	})
}
