package gox

import "net/http"

func (gx *Gox) GET(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return gx.Handle(pattern, http.HandlerFunc(handler), http.MethodGet)
}

func (gx *Gox) HEAD(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return gx.Handle(pattern, http.HandlerFunc(handler), http.MethodHead)
}

func (gx *Gox) POST(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return gx.Handle(pattern, http.HandlerFunc(handler), http.MethodPost)
}

func (gx *Gox) PUT(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return gx.Handle(pattern, http.HandlerFunc(handler), http.MethodPut)
}

func (gx *Gox) PATCH(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return gx.Handle(pattern, http.HandlerFunc(handler), http.MethodPatch)
}

func (gx *Gox) DELETE(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return gx.Handle(pattern, http.HandlerFunc(handler), http.MethodDelete)
}

func (gx *Gox) CONNECT(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return gx.Handle(pattern, http.HandlerFunc(handler), http.MethodConnect)
}

func (gx *Gox) OPTIONS(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return gx.Handle(pattern, http.HandlerFunc(handler), http.MethodOptions)
}

func (gx *Gox) TRACE(pattern string, handler func(w http.ResponseWriter, r *http.Request)) *routers {
	return gx.Handle(pattern, http.HandlerFunc(handler), http.MethodTrace)
}

func (gx *Gox) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request), methods ...string) *routers {
	return gx.Handle(pattern, http.HandlerFunc(handler), methods...)
}
