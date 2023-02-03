package main

import (
	"fmt"
	"net/http"

	"github.com/lightsaid/gotk/gox"
)

func main() {
	mux := gox.New()

	mux.GET("/get1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "get")
	})

	mux.POST("/post", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "post")
	})

	mux.DELETE("/delete", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "delete")
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
