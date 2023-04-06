package mux

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestNewTrie(t *testing.T) {
	tire := NewTrie()
	trieBytes, err := json.MarshalIndent(tire, "", " ")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(trieBytes))
}

func TestParse(t *testing.T) {
	segs := []struct {
		part  string
		key   string
		regex string
	}{
		{
			part:  "api",
			key:   "",
			regex: "",
		},
		{
			part:  ":name",
			key:   "name",
			regex: "",
		},
		{
			part:  ":id|^[0-9]$",
			key:   "id",
			regex: "^[0-9]$",
		},
		{
			part:  ":id|",
			key:   "id",
			regex: "",
		},
	}

	tree := NewTrie()
	for _, s := range segs {
		key, regex := tree.parse(s.part)
		if key != s.key {
			t.Errorf("%s key want: %s, got: %s", s.part, s.key, key)
		}
		if regex != s.regex {
			t.Errorf("%s regex want: %s, got: %s", s.part, s.regex, regex)
		}
	}
}

func TestInsertAndMatch(t *testing.T) {
	var testCase = []struct {
		pattern string
		err     error
		method  string
	}{
		{pattern: "/", err: nil, method: http.MethodGet},
		{pattern: "/api/:id", err: nil, method: http.MethodGet},
		{pattern: "/api/:id", err: nil, method: http.MethodPost},
		{pattern: "/api/:abc", err: ErrConflict, method: http.MethodGet},
		{pattern: "/api/:name", err: ErrConflict, method: http.MethodPost},
		{pattern: "/api/:name/age", err: nil, method: http.MethodDelete},
		{pattern: "/api/:title/:pid", err: nil, method: http.MethodPut},
		{pattern: "/api/:title/:pid|^[0-9]+$", err: nil, method: http.MethodOptions},
		{pattern: "/api/:title/:kid|^[A-Z]+$", err: nil, method: http.MethodHead},
		{pattern: "/api/:title/:kid|^[A-Z]+$/hello", err: nil, method: http.MethodPatch},

		{pattern: "/v1/auth/login", err: nil, method: http.MethodPost},
		{pattern: "/v1/auth/register", err: nil, method: http.MethodPost},

		{pattern: "/v1/:*", err: nil, method: http.MethodPost},
	}

	var tmpTrie = NewTrie()
	for _, test := range testCase {
		_, err := tmpTrie.Insert(test.pattern, nil, test.method)
		if err != nil && err != test.err && !errors.Is(err, ErrConflict) {
			t.Errorf("%s err want: %s, got: %s", test.pattern, test.err, err)
		}
	}
	// buf, _ := json.MarshalIndent(tmpTrie, "", " ")
	// fmt.Println(string(buf))

	var apis = []struct {
		path   string
		method string
		exists bool
	}{
		{
			path:   "/",
			method: http.MethodGet,
			exists: true,
		},
		{
			path:   "/",
			method: http.MethodPost,
			exists: false,
		},
		{
			path:   "/api/100",
			method: http.MethodGet,
			exists: true,
		},
		{
			path:   "/api/200",
			method: http.MethodPost,
			exists: true,
		},
		{
			path:   "/api/200/cat",
			method: http.MethodDelete,
			exists: false,
		},
		{
			path:   "/api/dog/222",
			method: http.MethodPut,
			exists: true,
		},
		{
			path:   "/api/dog/222",
			method: http.MethodPost,
			exists: false,
		},
		{
			path:   "/api/dog/333",
			method: http.MethodOptions,
			exists: true,
		},
		{
			path:   "/api/dog/MYSQL",
			method: http.MethodHead,
			exists: true,
		},
		{
			path:   "/api/dog/MYSQL/hello",
			method: http.MethodHead,
			exists: false,
		},
		{
			path:   "/api/dog/MYSQL/hello",
			method: http.MethodPatch,
			exists: true,
		},
		{path: "/v1/auth/login", exists: true, method: http.MethodPost},
		{path: "/v1/auth/register", exists: true, method: http.MethodPost},
		{path: "/v1/auth/login", exists: false, method: http.MethodGet},
		{path: "/v1/auth/register", exists: false, method: http.MethodPut},

		{path: "/v1/op/ff", exists: false, method: http.MethodPost},
	}

	for _, api := range apis {
		start := time.Now()
		request, err := http.NewRequest(api.method, api.path, nil)
		if err != nil {
			t.Error(err)
		}
		_, ok := tmpTrie.Match(request)
		if ok != api.exists {
			t.Errorf("%s %s want: %t, got: %t", api.method, api.path, api.exists, ok)
		}
		fmt.Println("duration: ", time.Since(start))
	}
}

/*
	如何处理路冲突和优先匹配问题？有可能出现如下情况:
	0和1冲突,不允许注册; 2和3不冲突,2优先匹配; 2~5都不冲突,2优先匹配,3~5随机优先匹配，寻找第一个符合条件的

	0. /api/:id
	1. /api/:name
	2. /api/:name/age
	3. /api/:title/:pid
	4. /api/:title/:pid|^[0-9]$
	5. /api/:title/:kid|^[A-C]$
	6. /api/:title/:kid|^[A-C]$/hello

	那么问题来了？假设有个接口: /api/cat/10，在寻找 /api/cat 这部分时，2～6 都有可能命中，按理说这个应该匹配到第3、4个，
	但是如果在匹配cat过程中优先命中2、5、6之一，则结果就是404了

	为了解决以上问题，因此在存储路由的时候，将类似 "/:id"、 "/:name"、 "/:title"... 这类动态路由结点key设为 * 结点;
	而 “/:pid|^[0-9]$”、"/:kid|^[A-C]$" 分别设置为 "^[0-9]$" 和 "^[A-C]$"
	最终上面 1~6 的存储结构大致如下:

				api

				:*       如何解决paramKe问题？(id、name、title)

	age    :*   :pid|^[0-9]$    :kid|^[A-C]$
									hello


	那么 * 结点要有足够的信息来做判断，因此在 Node 添加 regexes 字段来存正则，如果是没有正则匹配的动态路由，则设为默认正则 ".+"

	TODO: 路径清理, / "" 的路径

		// NOTE: path="/xx/" 转化为 path=“/xx”
	// NOTE: path="yyy" 转化为 path="/yyy"

*/
