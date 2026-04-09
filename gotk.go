package gotk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

var (
	baseMBSize       = 8
	ReadJSONMaxBytes = baseMBSize << 20 // 8MB
)

type Map map[string]any

type JSONResponse struct {
	BizCode   string `json:"bizCode"`             // 业务编码
	Message   string `json:"message"`             // 客户消息
	Data      any    `json:"data"`                // 任意数据
	Version   string `json:"version,omitempty"`   // 版本信息
	RequestID string `json:"requestId,omitempty"` // 请求Id，做简单的链路追踪
}

// ReadJSON 读取入参，绑定到 dst 上
func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// 限制请求体大小
	r.Body = http.MaxBytesReader(w, r.Body, int64(ReadJSONMaxBytes))

	// 使用请求体创建一个解码器
	dec := json.NewDecoder(r.Body)

	// dec.DisallowUnknownFields() // 不允许存非定义的字段，这里注释掉表示允许

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var typeError *json.UnmarshalTypeError
		var invalidError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return errors.New("请输入JSON格式参数")
		case errors.As(err, &typeError):
			return errors.New("参数字段类型不匹配")
		case errors.As(err, &invalidError):
			// 通常是dst传递的不是指针
			return errors.New("无法解码到目标参数")
		// case strings.Contains(err.Error(), "non-pointer"):
		// 	// 通常是dst传递的不是指针
		// 	return errors.New("无法解码到目标参数")
		case errors.Is(err, io.EOF):
			return errors.New("参数不能为空")
		case strings.Contains(err.Error(), "http: request body too large"):
			return fmt.Errorf("参数大小不能超过 %d MB", baseMBSize)

		default:
			slog.Error("dec.Decode fail", "error", err)
			return errors.New("未知错误，请检查参数")
		}
	}

	// 尝试在解码一次
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("请求体只能包含一个JSON值")
	}

	return nil
}

// WriteJSON 写入数据,如果a为nil，则默认请求成功
func WriteJSON(w http.ResponseWriter, r *http.Request, a *ApiError, data any, headers ...http.Header) error {
	if a == nil {
		// 使用保留编码
		a = newApiError(http.StatusOK, "-200", "请求成功")
	}

	// 组织返回数据
	response := &JSONResponse{
		BizCode:   a.BizCode(),
		Message:   a.Message(),
		Data:      data,
		Version:   GetByCtx(r, VersionCtxKey, ""),
		RequestID: GetByCtx(r, RequestIDCtxKey, ""),
	}

	out, err := json.Marshal(response)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for _, header := range headers {
			for key, values := range header {
				for _, v := range values {
					w.Header().Add(key, v)
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(a.StatusCode())
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}
