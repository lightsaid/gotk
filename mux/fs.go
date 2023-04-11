package mux

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
)

// 参考 gin 的实现方式 并简化

// Static 静态文件资源服务
func (s *ServeMux) Static(pattern string, dir string) {
	//1.  生成一个静态文件服务 handler,  http.Dir 实现了 http.FileSystem open 接口
	handler := s.createStaticHandler(pattern, http.Dir(dir))

	urlPattern := path.Join(pattern, "/:filepath")
	urlPattern = cleanPath(urlPattern)

	// 2. 注册静态文件服务
	s.Handle(urlPattern, handler, http.MethodGet, http.MethodHead)
}

// createStaticHandler 创建一个 http.Handler 给静态文件 handler
func (s *ServeMux) createStaticHandler(pattern string, fs http.FileSystem) http.Handler {
	handler := http.StripPrefix(pattern, http.FileServer(fs))
	// 直接返回可以访问整个目录，不安全，因此需要打开具体的文件
	// return handler

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filepatn := Param(r, "filepath")
		if filepatn == "" {
			log.Println("filepatn is nil")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		filename := strings.Replace(filepatn, pattern, "", 1)
		fmt.Println("filename: ", filepatn, "---", filename)
		// 尝试打开文件
		f, err := fs.Open(filename)
		if err != nil {
			log.Println("fs.Open failed: ", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer f.Close()

		if info, err := f.Stat(); err != nil {
			log.Println("get FileInfo failed: ", err)
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			if info.IsDir() {
				log.Println("禁止访问目录")
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}

		// 最终还是由  http.FileServer 返回文件
		handler.ServeHTTP(w, r)
	})
}
