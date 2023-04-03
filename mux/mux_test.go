package mux

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandle(t *testing.T) {
	var mux = NewServeMux()
	var handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(" Target ")
	})

	mux.Use(helloMiddleware)

	mux.Handle("/v1/auth/login", handlerFunc, http.MethodPost).Use(sinceMiddleware, loggerMiddleware)
	mux.Handle("/v2/auth/login", handlerFunc, http.MethodPut)

	// buf, _ := json.MarshalIndent(mux.routes.root, "", " ")
	// fmt.Println(string(buf))

	// mux.GET()

	request, err := http.NewRequest(http.MethodPost, "/v1/auth/login", nil)
	if err != nil {
		t.Error(err)
	}

	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)
}

func TestServeMux(t *testing.T) {
	var testCases = []struct {
		name       string
		url        string
		method     string
		statusCode int
		res        map[string]string
	}{
		{
			name:       "v1 Login Get Success",
			url:        "/v1/auth/login",
			method:     http.MethodGet,
			statusCode: 200,
			res:        map[string]string{"pattern": "/v1/auth/login"},
		},
		{
			name:       "v2 Login Delete MethodNotAllowed",
			url:        "/v2/auth/login",
			method:     http.MethodDelete,
			statusCode: 405,
			res:        map[string]string{},
		},

		{
			name:       "Api Posts",
			url:        "/api/posts",
			method:     http.MethodGet,
			statusCode: 200,
			res:        map[string]string{"pattern": "/api/posts"},
		},
		{
			name:       "Api Posts Title",
			url:        "/api/posts/heloo",
			method:     http.MethodGet,
			statusCode: 200,
			res:        map[string]string{"pattern": "/api/posts/:title", "title": "heloo"},
		},
		{
			name:       "Api Posts ID",
			url:        "/api/posts/100",
			method:     http.MethodPost,
			statusCode: 200,
			res:        map[string]string{"pattern": "/api/posts/:id|^[0-9]+$", "id": "100"},
		},

		{
			name:       "Api Posts Tag With Name 200",
			url:        "/api/posts/222/go",
			method:     http.MethodPost,
			statusCode: 200,
			res:        map[string]string{"pattern": "/api/posts/:tag|^[0-9]+$/:name", "tag": "222", "name": "go"},
		},
		{
			name:       "Api Posts Tag With Name 404",
			url:        "/api/posts/abc/go",
			method:     http.MethodPost,
			statusCode: 404,
			res:        map[string]string{},
		},
		{
			name:       "Api Posts Tag With Name 200",
			url:        "/api/posts/500/go/all",
			method:     http.MethodPost,
			statusCode: 200,
			res:        map[string]string{"pattern": "/api/posts/:tag|^[0-9]+$/:name/all", "tag": "500", "name": "go"},
		},
		// TODO 更多测试...
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(test.method, test.url, nil)
			if err != nil {
				t.Error(err)
			}
			response := httptest.NewRecorder()
			testMux.ServeHTTP(response, request)

			if request.Method != test.method {
				t.Errorf("url: %s, method want: %s, got: %s", test.url, test.method, request.Method)
			}

			if test.statusCode != response.Code {
				t.Errorf("url: %s,  statusCode want: %d, got: %d", test.url, test.statusCode, response.Code)
			} else if test.statusCode == http.StatusOK && response.Code == http.StatusOK {
				var data map[string]string
				err := json.NewEncoder(response).Encode(&data)
				if err != nil {
					t.Error(err)
				}
				if v, ok := test.res["id"]; ok {
					fmt.Println("v: ", v, data["id"])
					if v != data["id"] {
						t.Errorf("url: %s, id want: %s, got: %s", test.url, test.res["id"], v)
					}
				}
			}
		})
	}
}
