package mux

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

var testMux *ServeMux

func TestMain(m *testing.M) {
	var patterns = []string{
		"/v1/auth/login",
		"/v1/auth/register",
		"/v2/auth/login",
		"/v2/auth/register",
		"/v1/auth/login/admin",
		"/v2/auth/login/admin",

		"/api/posts",
		"/api/posts/:title",
		"/api/posts/:id|^[0-9]+$",
		"/api/posts/:tag|^[0-9]+$/:name",
		"/api/posts/:tag|^[0-9]+$/:name/all",

		"/api/tags",
		"/api/tags/:name",
		"/tags/:name",
	}

	testMux = NewServeMux()

	for _, pattern := range patterns {
		testMux.Handle(pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := map[string]string{
				"pattern": pattern,
				"title":   Param(r, "title"),
				"id":      Param(r, "id"),
				"tag":     Param(r, "tag"),
				"name":    Param(r, "name"),
			}
			_ = json.NewEncoder(w).Encode(data)
		}), http.MethodGet, http.MethodPost) // 这里仅仅接收 Get 和 Post 请求
	}

	os.Exit(m.Run())
}
