package gox

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var routeRxMaps = make(map[string]regexp.Regexp)

type contextKey string

type route struct {
	h       http.Handler
	paths   []string
	methods []string
	pattern string
}

func Param(r *http.Request, key string) string {
	val, ok := r.Context().Value(contextKey(key)).(string)
	if !ok {
		return ""
	}
	return val
}

func (r *route) match(ctx context.Context, paths []string) (context.Context, bool) {
	fmt.Println("mathc >>>", r.paths, paths, len(r.paths) != len(paths), r.pattern)
	if len(r.paths) != len(paths) {
		return ctx, false
	}

	rex := regexp.MustCompile(`{.+}`)

	var exist bool
	for i, path := range paths {
		// 查找是否存在 {.+} 部分内容
		parts := rex.FindStringSubmatch(r.paths[i])
		fmt.Println("parts 1: ", parts, len(parts), path, r.paths[i])
		// 不存在 {} parmas, 并且路由不匹配
		if len(parts) == 0 && r.paths[i] != path {
			fmt.Println("aaa")
			return ctx, false
		} else if len(parts) == 0 && r.paths[i] == path {
			fmt.Println("bbb", i, len(paths)-1)
			// 这一部分匹配成功，标记
			exist = true
			continue
		}

		fmt.Println("exists: ", true)

		// parts 存在 {} 部分， 提取；如果有正则表达式则分割出来

		// 仅仅只有{},里面是空的
		if len(parts[0]) == 2 {
			return ctx, false
		}

		// 提取, 如： {id:[0-9]+} => id:[0-9]+
		part := parts[0][1 : len(parts[0])-1]

		if strings.Contains(part, ":") {
			// part 包含参数名和正则
			ps := strings.Split(part, ":")
			if len(ps) != 2 {
				return ctx, false
			}
			rx, ok := routeRxMaps[ps[1]]
			if !ok {
				// 不存在则添加到缓存里
				routeRxMaps[ps[1]] = *regexp.MustCompile(ps[1])
				rx = routeRxMaps[ps[1]]
			}
			if exist = rx.MatchString(path); exist {
				fmt.Println(">>> 1 param: ", part, path)
				ctx = context.WithValue(ctx, contextKey(ps[0]), path)
			} else {
				return ctx, false
			}
		} else {
			// part 就是 parmas 的参数名
			fmt.Println(">>> 2 param: ", part, path)
			ctx = context.WithValue(ctx, contextKey(part), path)
		}
	}

	return ctx, exist
}
