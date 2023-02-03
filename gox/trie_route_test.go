package gox_test

import (
	"context"
	"strings"
	"testing"

	"github.com/lightsaid/gotk/gox"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name     string
	route    string
	expected bool
}

func registerRoute(t *testing.T) gox.RouteTrie {
	tests := []testCase{
		{name: "v1正常路由", route: "/v1/api/users", expected: true},
		{name: "v2正常路由", route: "/v2/api/users", expected: true},
		{name: "v2重复路由", route: "/v2/api/users", expected: false},
		{name: "v1 id 参数路由", route: "/v1/api/users/:id|[0-9]+", expected: true},
		{name: "v1", route: "/v2/api/users/:id|[0-9]+/mine", expected: true},
		{name: "v1", route: "/v1/api/users/sys/:id|[0-9]+", expected: true},
		{name: "v1", route: "/v2/api/users/sys/:id|[0-9]+/profile", expected: true},
		{name: "v1", route: "/v2/api/users/sys/:id|[0-9]+/profile", expected: false},
		{name: "v1", route: "/v1/api/users/sys/:status|^[A-Z]+$", expected: true},
		{name: "v1", route: "/v1/api/users/sys/:status|^[A-Z]+$", expected: false},
		{name: "v1", route: "/v2/api/users/sys/:status|^[A-Z]+$", expected: true},
		{name: "v1", route: "/v1/api/users/:gender|^[A-Z]$/:age|[0-9]{1,3}$", expected: true},
		{name: "v1", route: "/v2/api/users/:gender|^[A-Z]$/:age|[0-9]{1,3}$/:name", expected: true},
		{name: "v1", route: "/api/products", expected: true},
		{name: "v1", route: "/api/category/list", expected: true},
		{name: "v1", route: "/api/category/get", expected: true},
		{name: "v1", route: "/api/category/post", expected: true},
		{name: "v1", route: "/api/category/post", expected: false},
		{name: "v1", route: "/api/category/put", expected: true},
		{name: "v1", route: "/api/category/del", expected: true},
		{name: "v1", route: "/api/category/del", expected: false},
		{name: "v1", route: "/api/products/:category/:id", expected: true},
	}

	trie := gox.NewRouteTrie()
	count := 0
	for _, test := range tests {
		parts := strings.Split(test.route, "/")[1:]
		ok := trie.Register(parts)
		require.Equal(t, test.expected, ok)
		if test.expected {
			count++
		}
	}
	require.Equal(t, trie.Size(), count)

	return trie
}

func TestTrieRegister(t *testing.T) {
	_ = registerRoute(t)
}

func TestParams(t *testing.T) {

	type testCase struct {
		path       string
		paramKey   string
		paramValue string
		match      bool
		route      string
	}

	tests := []testCase{
		{path: "/v1/api/users", paramKey: "", paramValue: "", match: true, route: "/v1/api/users"},
		{path: "/v2/api/users", paramKey: "", paramValue: "", match: true, route: "/v2/api/users"},
		{path: "/v1/api/users/100", paramKey: "id", paramValue: "100", match: true, route: "/v1/api/users/:id|[0-9]+"},
		{path: "/v2/api/users/200/mine", paramKey: "id", paramValue: "200", match: true, route: "/v2/api/users/:id|[0-9]+/mine"},
		{path: "/v1/api/users/abc", paramKey: "id", paramValue: "abc", match: false, route: "/v1/api/users/:id|[0-9]+"},
		{path: "/v2/api/users/300/mine", paramKey: "id", paramValue: "300", match: true, route: "/v2/api/users/:id|[0-9]+/mine"},
		{path: "/v1/api/users/sys/400", paramKey: "id", paramValue: "400", match: true, route: "/v1/api/users/sys/:id|[0-9]+"},
		{path: "/v2/api/users/sys/500/profile", paramKey: "id", paramValue: "500", match: true, route: "/v2/api/users/sys/:id|[0-9]+/profile"},
		{path: "/v2/api/users/sys/abc/profile", match: false, route: "/v2/api/users/sys/:id|[0-9]+/profile"},
		{path: "/v1/api/users/sys/STOP", paramKey: "status", paramValue: "STOP", match: true, route: "/v1/api/users/sys/:status|^[A-Z]+$"},
		{path: "/v1/api/users/sys/stop", match: false, route: "/v1/api/users/sys/:status|^[A-Z]+$"},
		{path: "/v2/api/users/sys/stop", match: false, route: "/v2/api/users/sys/:status|^[A-Z]+$"},
		{path: "/v2/api/users/sys/ABC", paramKey: "status", paramValue: "ABC", match: true, route: "/v2/api/users/sys/:status|^[A-Z]+$"},
		{path: "/v1/api/users/F/88", paramKey: "gender", paramValue: "F", match: true, route: "/v1/api/users/:gender|^[A-Z]$/:age|[0-9]{1,3}$"},
		{path: "/v1/api/users/F/88", paramKey: "age", paramValue: "88", match: true, route: "/v1/api/users/:gender|^[A-Z]$/:age|[0-9]{1,3}$"},
		{path: "/v1/api/users/F/aa", match: false, route: "/v1/api/users/:gender|^[A-Z]$/:age|[0-9]{1,3}$"},
		{path: "/v1/api/users/m/99", match: false, route: "/v1/api/users/:gender|^[A-Z]$/:age|[0-9]{1,3}$"},
		{path: "/v2/api/users/M/88/mario", paramKey: "gender", paramValue: "M", match: true, route: "/v2/api/users/:gender|^[A-Z]$/:age|[0-9]{1,3}$/:name"},
		{path: "/v2/api/users/M/88/mario", paramKey: "age", paramValue: "88", match: true, route: "/v2/api/users/:gender|^[A-Z]$/:age|[0-9]{1,3}$/:name"},
		{path: "/v2/api/users/M/88/mario", paramKey: "name", paramValue: "mario", match: true, route: "/v2/api/users/:gender|^[A-Z]$/:age|[0-9]{1,3}$/:name"},
		{path: "/v2/api/users/MALE/88/mario", match: false, route: "/v2/api/users/:gender|^[A-Z]$/:age|[0-9]{1,3}$/:name"},
		{path: "/api/category/del", match: true, route: "/api/category/del"},
		{path: "/api/products/phone/111", match: true, paramKey: "category", paramValue: "phone", route: "/api/products/:category/:id"},
		{path: "/api/products/phone/111", match: true, paramKey: "id", paramValue: "111", route: "/api/products/:category/:id"},
	}

	trie := registerRoute(t)
	for _, test := range tests {
		ctx := context.Background()
		parts := strings.Split(test.path, "/")[1:]
		// fmt.Println(">> parts: ", parts)
		c, ok := trie.Match(ctx, parts)
		require.Equal(t, test.match, ok)
		if len(test.paramKey) > 0 && test.match {
			val := c.Value(gox.ContextKey(test.paramKey))
			require.Equal(t, test.paramValue, val)
		}

		var matchCount = 0
		var keyCount = 0
		if test.match {
			matchCount++
			if key, ok := c.Value(gox.ContextKey("trieNodeKey")).(string); ok {
				require.Equal(t, test.route, key)
				keyCount++
			}
		}

		require.Equal(t, matchCount, keyCount)
	}
}
