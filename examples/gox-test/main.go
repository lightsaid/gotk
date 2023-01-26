package main

import (
	"fmt"
	"net/http"

	"github.com/lightsaid/gotk/gox"
)

func main() {
	mux := gox.New()

	mux.GET("/get", func(w http.ResponseWriter, r *http.Request) {
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

	http.ListenAndServe("0.0.0.0:9999", mux)
}
