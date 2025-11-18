package gotk

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type CtxKey string

var (
	VersionCtxKey   CtxKey = "gotk_version"
	RequestIDCtxKey CtxKey = "gotk_request_id"
)

// SetApiVersion 设置版本信息到上下文
func SetVersionCtx(next http.Handler, version string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), VersionCtxKey, version)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// SetRequestIDCtx 设置request id
func SetRequestIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.NewString()
		ctx := context.WithValue(r.Context(), RequestIDCtxKey, requestId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetByCtx
func GetByCtx[T any](r *http.Request, key CtxKey, defaultVal T) T {
	val, exist := r.Context().Value(key).(T)
	if !exist {
		return defaultVal
	}
	return val
}
