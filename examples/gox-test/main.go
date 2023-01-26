package main

import (
	"fmt"
	"net/http"

	"github.com/lightsaid/gotk/gox"
)

func main() {
	mux := gox.New()
	mux.HandleFunc("/", indexHandle)
	mux.HandleFunc("/tag", tagHandle)
	mux.HandleFunc("/post", postHandle)

	http.ListenAndServe("0.0.0.0:9999", mux)

}

func indexHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "indexHandle")
}

func tagHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "tagHandle")
}

func postHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "postHandle")
}
