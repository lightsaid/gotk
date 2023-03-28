package mux

import (
	"context"
	"strings"
	"testing"
)

func TestTrieInsertAndMatch(t *testing.T) {
	testCases := []struct {
		pattern string
		url     string
		exists  bool
	}{
		{pattern: "/v1/auth/login", url: "/v1/auth/login", exists: true},
		{pattern: "/v2/auth/login", url: "/v2/auth/login", exists: true},

		{pattern: "", url: "/v3/auth/login", exists: false},
		{pattern: "", url: "/v2/auth", exists: false},
		{pattern: "", url: "/v2/auth/login/test", exists: false},

		{pattern: "/v1/auth/register", url: "/v1/auth/register", exists: true},
		{pattern: "/v1/auth1/register", url: "/v1/auth1/register", exists: true},
		{pattern: "", url: "v1/auth1", exists: false},
		{pattern: "", url: "v1/auth1/registerss", exists: false},

		{pattern: "/v1/api/product/:title", url: "/v1/api/product/macbook", exists: true},
		{pattern: "", url: "/v1/api/product/huawei", exists: true},
		{pattern: "", url: "/v1/api/product/100", exists: true},

		{pattern: "/v1/api/products/:id|^[0-9]+$", url: "/v1/api/products/100", exists: true},
		{pattern: "", url: "/v1/api/products/abc", exists: false},

		{pattern: "/v1/api/products/:id|^[0-9]+$/:name", url: "/v1/api/products/100/abc", exists: true},
		{pattern: "/v2/api/products/:id|^[0-9]+$/:name", url: "/v2/api/products/100/abc", exists: true},
		{pattern: "", url: "/v2/api/products/abc", exists: false},
	}

	testTrie := NewTrie()
	for _, test := range testCases {
		if len(test.pattern) > 0 {
			testTrie.Insert(strings.Split(test.pattern, "/")[1:])
		}
	}

	for _, test := range testCases {
		_, ok := testTrie.Match(context.Background(), strings.Split(test.url, "/")[1:])
		if ok != test.exists {
			t.Errorf("url: %s want: %t, got: %t", test.url, test.exists, ok)
		}
	}
}
