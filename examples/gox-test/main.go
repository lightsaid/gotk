package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

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

	mux.GET("/api/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		res := gox.Param(r, "id")
		fmt.Fprintf(w, "param: %s", res)
	})

	mux.GET("/api/tag/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		res := gox.Param(r, "id")
		fmt.Fprintf(w, "param: %s", res)
	})

	http.ListenAndServe("0.0.0.0:9999", mux)

	// var s = "/products/{key}"
	// var s = "/articles/{category}/{id}"
	// var s = "/articles/{category}/{id:[0-9]+}"
	var s = "/articles/{id:[0-9]+}"

	// var r = "/articles/tag/90"

	var ss = strings.Split(s, "/")
	fmt.Println(ss, len(ss))

	for i := 0; i < len(ss); i++ {
		if ss[i] == "" {
			fmt.Println(i, "nil")
		}
		fmt.Printf("[%d]=%q\n", i, ss[i])
	}

	// paths := strings.Split(r, "/")

	// for _, path := range paths {

	// }

	// 提取正则表达式
	rex := regexp.MustCompile(`{.+}`)

	subs := rex.FindStringSubmatch(s)

	part := subs[0][1 : len(subs[0])-1]
	fmt.Println(subs, len(subs), part)

	ps := strings.Split(part, ":")
	fmt.Println(ps, len(ps))

}
