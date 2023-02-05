package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/lightsaid/gotk/gox"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("%s - %s %s %s\n", r.RemoteAddr, r.Proto, r.Method, r.URL)

		w.Header().Set("abc", "abc")
		ctx := context.WithValue(r.Context(), gox.ContextKey("fff"), "hello fff")
		fmt.Println("<<<<< AAA >>>>>")

		// rad := rand.Intn(10)
		// if rad > 5 {
		// 	next.ServeHTTP(w, r.WithContext(ctx))
		// }
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func bbbRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("<<<<< BBB >>>>>")
		fmt.Println("bbb result: ", w.Header().Get("abc"), r.Context().Value(gox.ContextKey("fff")))
		next.ServeHTTP(w, r)
	})
}

func cccRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("<<<<< CCC >>>>>")
		fmt.Println("ccc result: ", w.Header().Get("abc"), r.Context().Value(gox.ContextKey("fff")))
		next.ServeHTTP(w, r)
	})
}

func dddRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("<<<<< DDD >>>>>")
		fmt.Println("ddd result: ", w.Header().Get("abc"), r.Context().Value(gox.ContextKey("fff")))
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := gox.New()

	mux.Register("AAA", logRequest)
	mux.Register("BBB", bbbRequest)
	mux.Register("CCC", cccRequest)
	mux.Register("DDD", dddRequest)

	mux.Global("AAA", "BBB")

	mux.GET("/get", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("result: ", w.Header().Get("abc"), r.Context().Value(gox.ContextKey("fff")))
		fmt.Fprintf(w, "get")
		fmt.Println("完成请求～～")
	})

	mux.POST("/post", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "post")
		fmt.Println("post 请求完成～～")
	}).Before("CCC").After("DDD")

	mux.DELETE("/delete", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "delete")
		fmt.Println("post 完成请求～～")
	})

	mux.PUT("/put", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "put")
	})

	mux.GET("/api/:id|[0-9]+", func(w http.ResponseWriter, r *http.Request) {
		res := gox.Param(r, "id")
		fmt.Fprintf(w, "param: %s", res)
	})

	mux.GET("/api/tag/:id|[0-9]+", func(w http.ResponseWriter, r *http.Request) {
		res := gox.Param(r, "id")
		fmt.Fprintf(w, "get request: param: %s", res)
	})

	mux.POST("/api/tag/:id|[0-9]+", func(w http.ResponseWriter, r *http.Request) {
		res := gox.Param(r, "id")
		fmt.Fprintf(w, "post request: param: %s", res)
	})

	http.ListenAndServe("0.0.0.0:9999", mux)

}
