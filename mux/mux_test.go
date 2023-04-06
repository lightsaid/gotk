package mux_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lightsaid/gotk/mux"
)

const (
	version         = "0.01"
	hello_mw_output = ">>> Hello Middleware <<<"
	since_mw_output = ">>> Execute Time:"
)

type TestTable struct {
	name   string
	url    string
	method string
	params *map[string]string
	status int
	result string
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", version)
}

func handlerEchoURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", r.URL.Path)
}

func myHelloMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(hello_mw_output)
		handler.ServeHTTP(w, r)
	})
}

func mySinceMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf("%s %s %s %s \n", r.Method, r.URL.Path, since_mw_output, time.Since(t))

	})
}

// 注册测试所需的路由
func registerRoutes() {
	// Version Handler
	testMux.Handle("/api/version", http.HandlerFunc(getVersion), http.MethodGet, http.MethodPost)
	testMux.HandleFunc("/api/v2/version", getVersion, http.MethodGet, http.MethodPost)

	// GetByID
	testMux.GET("/api/:id|^[0-9]+$", handlerEchoURL)

	// NormalRoute 不太合理的api测试
	testMux.GET("/api/:name", handlerEchoURL)
	testMux.GET("/api/:name/:id|^[0-9]+", handlerEchoURL)
	testMux.GET("/api/:cat/cat", handlerEchoURL)
	testMux.GET("/api/:dog/dog", handlerEchoURL)
	testMux.GET("/api/:fish/fish", handlerEchoURL)
	testMux.GET("/api/:fish/fish/:id|^[a-zA-Z]+$", handlerEchoURL)
}

func TestVersionHandler(t *testing.T) {
	registerRoutes()

	var testCase = []TestTable{
		{url: "/api/version", method: http.MethodGet, status: 200, result: version},
		{url: "/api/version", method: http.MethodPost, status: 200, result: version},
		{url: "/api/version", method: http.MethodPut, status: 404, result: ""},

		{url: "/api/v2/version", method: http.MethodPost, status: 200, result: version},
		{url: "/api/v2/version", method: http.MethodDelete, status: 404, result: ""},
	}

	for _, tc := range testCase {
		t.Run("Version Handler", func(t *testing.T) {
			request, err := http.NewRequest(tc.method, tc.url, nil)
			if err != nil {
				t.Error(err)
			}
			response := httptest.NewRecorder()

			testMux.ServeHTTP(response, request)

			result := response.Body.String()

			if tc.method != request.Method {
				t.Errorf("%s %s method want: %s, got: %s", tc.method, tc.url, tc.method, request.Method)
			}

			if tc.status != response.Code {
				t.Errorf("%s %s statusCode want: %d, got: %d", tc.method, tc.url, tc.status, response.Code)
			}

			if tc.status == 200 && version != result {
				t.Errorf("%s %s version want: %s, got: %s", tc.method, tc.url, tc.result, result)
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	registerRoutes()

	var testCase = []TestTable{
		{url: "/api/version", method: http.MethodGet, status: 200, result: version},
		{url: "/api/version", method: http.MethodPut, status: 405, result: ""},
	}

	testMux.OpenAllowed()

	for _, tc := range testCase {
		request, err := http.NewRequest(tc.method, tc.url, nil)
		if err != nil {
			t.Error(err)
		}
		response := httptest.NewRecorder()
		testMux.ServeHTTP(response, request)

		if tc.status != response.Code {
			t.Errorf("%s %s statusCode want: %d, got: %d", tc.method, tc.url, tc.status, response.Code)
		}
	}
}

func TestGetByID(t *testing.T) {
	registerRoutes()

	var testCase = []TestTable{
		{url: "/api/100", method: http.MethodGet, status: 200, result: "/api/100"},
		{url: "/api/200", method: http.MethodPost, status: 404, result: ""},
	}

	t.Run("GetByID", func(t *testing.T) {
		for _, tc := range testCase {
			request, err := http.NewRequest(tc.method, tc.url, nil)
			if err != nil {
				t.Error(err)
			}
			response := httptest.NewRecorder()

			testMux.ServeHTTP(response, request)

			result := response.Body.String()

			if tc.method != request.Method {
				t.Errorf("%s %s method want: %s, got: %s", tc.method, tc.url, tc.method, request.Method)
			}

			if tc.status != response.Code {
				t.Errorf("%s %s statusCode want: %d, got: %d", tc.method, tc.url, tc.status, response.Code)
			}

			if tc.status == 200 && tc.url != result {
				t.Errorf("%s %s result want: %q, got: %q", tc.method, tc.url, tc.url, result)
			}
		}
	})
}

func TestNormalRoutes(t *testing.T) {
	registerRoutes()

	mux.SetMuxMode("debug")
	testMux.PrintTrieRoutes()

	var testCase = []TestTable{
		{url: "/api/zhansan", method: http.MethodGet, status: 200},
		{url: "/api/zhansan/200", method: http.MethodGet, status: 200},
		{url: "/api/cat_001/cat", method: http.MethodGet, status: 200},
		{url: "/api/dog_001/dog", method: http.MethodGet, status: 200},
		{url: "/api/fish_001/fish", method: http.MethodGet, status: 200},
		{url: "/api/fish_001/fish/999", method: http.MethodGet, status: 404},
		{url: "/api/fish_001/fish/abc", method: http.MethodGet, status: 200},
	}

	t.Run("NormalRoute", func(t *testing.T) {
		for _, tc := range testCase {
			request, err := http.NewRequest(tc.method, tc.url, nil)
			if err != nil {
				t.Error(err)
			}
			response := httptest.NewRecorder()

			testMux.ServeHTTP(response, request)

			result := response.Body.String()

			if tc.method != request.Method {
				t.Errorf("%s %s method want: %s, got: %s", tc.method, tc.url, tc.method, request.Method)
			}

			if tc.status != response.Code {
				t.Errorf("%s %s statusCode want: %d, got: %d", tc.method, tc.url, tc.status, response.Code)
			}

			if tc.status == 200 && tc.url != result {
				t.Errorf("%s %s result want: %q, got: %q", tc.method, tc.url, tc.url, result)
			}
		}
	})
}

func TestGlobalMiddleware(t *testing.T) {
	registerRoutes()

	testMux.Use(myHelloMiddleware, mySinceMiddleware)

	var testCase = []TestTable{
		{url: "/api/lightsaid", method: http.MethodGet, status: 200},
		{url: "/api/kitty/cat", method: http.MethodGet, status: 200},
	}

	for _, tc := range testCase {

		// // 暂存一下 os.Stdout
		// saveStd := os.Stdout

		// r, w, err := os.Pipe()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// // 接收/代理终端输出
		// os.Stdout = w

		// 包装请求
		request, err := http.NewRequest(tc.method, tc.url, nil)
		if err != nil {
			t.Error(err)
		}
		response := httptest.NewRecorder()
		testMux.ServeHTTP(response, request)

		// // 关闭输出
		// _ = w.Close()
		// fmt.Println("close>>>>>")

		// // 获取终端输出
		// result, err := io.ReadAll(r)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// output := string(result)
		// os.Stdout = saveStd

		// fmt.Println(">>> output: ", output)
		// if !strings.Contains(output, hello_mw_output) {
		// 	t.Error("没有 hello_mw_output 输出")
		// }

		// if !strings.Contains(output, since_mw_output) {
		// 	t.Error("没有 since_mw_output 输出")
		// }
	}
}

// func TestHandle(t *testing.T) {
// 	var mux = NewServeMux()
// 	var handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(" Target ")
// 	})

// 	mux.Use(helloMiddleware)

// 	mux.Handle("/v1/auth/login", handlerFunc, http.MethodPost).Use(sinceMiddleware, loggerMiddleware)
// 	mux.Handle("/v2/auth/login", handlerFunc, http.MethodPut)

// 	// buf, _ := json.MarshalIndent(mux.routes.root, "", " ")
// 	// fmt.Println(string(buf))

// 	// mux.GET()

// 	request, err := http.NewRequest(http.MethodPost, "/v1/auth/login", nil)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	response := httptest.NewRecorder()

// 	mux.ServeHTTP(response, request)
// }

// func TestServeMux(t *testing.T) {
// 	var testCases = []struct {
// 		name       string
// 		url        string
// 		method     string
// 		statusCode int
// 		res        map[string]string
// 	}{
// 		{
// 			name:       "v1 Login Get Success",
// 			url:        "/v1/auth/login",
// 			method:     http.MethodGet,
// 			statusCode: 200,
// 			res:        map[string]string{"pattern": "/v1/auth/login"},
// 		},
// 		{
// 			name:       "v2 Login Delete MethodNotAllowed",
// 			url:        "/v2/auth/login",
// 			method:     http.MethodDelete,
// 			statusCode: 405,
// 			res:        map[string]string{},
// 		},

// 		{
// 			name:       "Api Posts",
// 			url:        "/api/posts",
// 			method:     http.MethodGet,
// 			statusCode: 200,
// 			res:        map[string]string{"pattern": "/api/posts"},
// 		},
// 		{
// 			name:       "Api Posts Title",
// 			url:        "/api/posts/heloo",
// 			method:     http.MethodGet,
// 			statusCode: 200,
// 			res:        map[string]string{"pattern": "/api/posts/:title", "title": "heloo"},
// 		},
// 		{
// 			name:       "Api Posts ID",
// 			url:        "/api/posts/100",
// 			method:     http.MethodPost,
// 			statusCode: 200,
// 			res:        map[string]string{"pattern": "/api/posts/:id|^[0-9]+$", "id": "100"},
// 		},

// 		{
// 			name:       "Api Posts Tag With Name 200",
// 			url:        "/api/posts/222/go",
// 			method:     http.MethodPost,
// 			statusCode: 200,
// 			res:        map[string]string{"pattern": "/api/posts/:tag|^[0-9]+$/:name", "tag": "222", "name": "go"},
// 		},
// 		{
// 			name:       "Api Posts Tag With Name 404",
// 			url:        "/api/posts/abc/go",
// 			method:     http.MethodPost,
// 			statusCode: 404,
// 			res:        map[string]string{},
// 		},
// 		{
// 			name:       "Api Posts Tag With Name 200",
// 			url:        "/api/posts/500/go/all",
// 			method:     http.MethodPost,
// 			statusCode: 200,
// 			res:        map[string]string{"pattern": "/api/posts/:tag|^[0-9]+$/:name/all", "tag": "500", "name": "go"},
// 		},
// 		// TODO 更多测试...
// 	}

// 	for _, test := range testCases {
// 		t.Run(test.name, func(t *testing.T) {
// 			request, err := http.NewRequest(test.method, test.url, nil)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			response := httptest.NewRecorder()
// 			testMux.ServeHTTP(response, request)

// 			if request.Method != test.method {
// 				t.Errorf("url: %s, method want: %s, got: %s", test.url, test.method, request.Method)
// 			}

// 			if test.statusCode != response.Code {
// 				t.Errorf("url: %s,  statusCode want: %d, got: %d", test.url, test.statusCode, response.Code)
// 			} else if test.statusCode == http.StatusOK && response.Code == http.StatusOK {
// 				var data map[string]string
// 				err := json.NewEncoder(response).Encode(&data)
// 				if err != nil {
// 					t.Error(err)
// 				}
// 				if v, ok := test.res["id"]; ok {
// 					fmt.Println("v: ", v, data["id"])
// 					if v != data["id"] {
// 						t.Errorf("url: %s, id want: %s, got: %s", test.url, test.res["id"], v)
// 					}
// 				}
// 			}
// 		})
// 	}
// }

// func TestNormalRoute(t *testing.T) {
// 	var testCase = []struct {
// 		pattern string
// 		err     error
// 		method  string
// 	}{
// 		{pattern: "/v1/:cat/cat", err: nil, method: http.MethodGet},
// 		{pattern: "/v1/:dog/dog", err: nil, method: http.MethodGet},
// 		{pattern: "/v1/:fish/fish", err: nil, method: http.MethodGet},
// 	}

// 	var tmpTrie = NewTrie()
// 	for _, test := range testCase {
// 		_, err := tmpTrie.Insert(test.pattern, nil, test.method)
// 		if err != nil && err != test.err && !errors.Is(err, ErrConflict) {
// 			t.Errorf("%s err want: %s, got: %s", test.pattern, test.err, err)
// 		}
// 	}

// 	buf, _ := json.MarshalIndent(tmpTrie, "", " ")
// 	fmt.Println(string(buf))

// 	var apis = []struct {
// 		path   string
// 		method string
// 		exists bool
// 	}{
// 		{
// 			path:   "/v1/cat001/cat",
// 			method: http.MethodGet,
// 			exists: true,
// 		},
// 		{
// 			path:   "/v1/dog002/dog",
// 			method: http.MethodGet,
// 			exists: true,
// 		},
// 		{
// 			path:   "/v1/fish003/fish",
// 			method: http.MethodGet,
// 			exists: true,
// 		},
// 	}

// 	for _, api := range apis {
// 		start := time.Now()
// 		request, err := http.NewRequest(api.method, api.path, nil)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		_, ok := tmpTrie.Match(request)
// 		if ok != api.exists {
// 			t.Errorf("%s %s want: %t, got: %t", api.method, api.path, api.exists, ok)
// 		}
// 		fmt.Println("duration: ", time.Since(start))
// 	}
// }
