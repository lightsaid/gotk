package gotk

import (
	"errors"
	"fmt"
	"strconv"
)

// ApiError 错误结构体，统一返回错误信息
type ApiError struct {
	message    string // 可读错误信息
	bizCode    string // 自定义唯一业务状态码
	statusCode int    // HTTP status code
	err        error  // 原始错误信息
}

// 收集自定义业务状态码，保持每个状态码都是唯一的
// 为了扩展性，自身保留[-1000, -100]的key
var (
	bizCodeMap = map[string]struct{}{}
	minCode    = -1000
	maxCode    = -100
)

// NewApiError 创建一个ApiError，bizCode 重复会触发panic;
// bizCode 是自定义业务编码;
// msg 是人类可读错误提示;
// statusCode 是 HTTP status code;
func NewApiError(statusCode int, bizCode, msg string) *ApiError {

	intCode, _ := strconv.Atoi(bizCode)
	if intCode >= minCode && intCode <= -100 {
		panic(fmt.Sprintf("%s is the reserved code", bizCode))
	}

	if _, exists := bizCodeMap[bizCode]; exists {
		panic(fmt.Sprintf("bizCode=%s already exist, please replace it", bizCode))
	}
	bizCodeMap[bizCode] = struct{}{}

	return newApiError(statusCode, bizCode, msg)
}

func newApiError(statusCode int, bizCode, msg string) *ApiError {
	return &ApiError{
		message:    msg,
		bizCode:    bizCode,
		statusCode: statusCode,
	}
}

// Error 实现 error 类型接口
func (a *ApiError) Error() string {
	return fmt.Sprintf(
		"statusCode=%d, bizCode=%s, message=%s, err=%v",
		a.statusCode,
		a.bizCode,
		a.message,
		a.err,
	)
}

// Unwrap 解开，提供给 errors.Is 和 errors.As 使用
func (a *ApiError) Unwrap() error {
	return a.err
}

// BizCode 返回自定义业务错误码
func (a *ApiError) BizCode() string {
	return a.bizCode
}

// Message 返回可读错误信息
func (a *ApiError) Message() string {
	return a.message
}

// StatusCode 返回HTTP Status Code
func (a *ApiError) StatusCode() int {
	return a.statusCode
}

// WithMessage 修改消息，返回一个新的 ApiError 指针
func (a *ApiError) WithMessage(msg string) *ApiError {
	e := *a
	e.message = msg
	return &e
}

// WithError 添加/追加错误, 返回一个新的 ApiError 指针
func (a *ApiError) WithError(err error) *ApiError {
	e := *a
	e.err = errors.Join(e.err, err)
	return &e
}

// With 添加error和message
func (a *ApiError) With(err error, msg string) *ApiError {
	return a.WithError(err).WithMessage(msg)
}

// 例子
// var (
// 	ErrNotFound = NewApiError(http.StatusNotFound, "NOT_FOUND", "查询不到")
// )

// func writeJSON(w http.ResponseWriter, r *http.Request, err *ApiError) {
// 	response := Map{"bizCode": err.BizCode(), "message": err.Message()}
// 	json.NewEncoder(w).Encode(response)
// }
