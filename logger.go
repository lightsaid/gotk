package gotk

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

type logType int

const (
	JsonType logType = iota + 1
	TextType
)

type LogHandleFunc func(ctx context.Context, r slog.Record) error

type LogHandler struct {
	slog.Handler
	handlers []LogHandleFunc
}

// Register 注册自定义Handle函数
func (lh *LogHandler) Register(handlers ...LogHandleFunc) {
	lh.handlers = append(lh.handlers, handlers...)
}

func (lh *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	requestID, ok := ctx.Value(RequestIDCtxKey).(string)
	if ok && requestID != "" {
		r.AddAttrs(slog.String("request_id", requestID))
	}

	version, ok := ctx.Value(VersionCtxKey).(string)
	if ok && version != "" {
		r.AddAttrs(slog.String("version", version))
	}

	var errs []error
	for _, fn := range lh.handlers {
		if err := fn(ctx, r); err != nil {
			errs = append(errs, err)
		}
	}

	err := lh.Handler.Handle(ctx, r)

	return errors.Join(err, errors.Join(errs...))
}

// NewLogger 创建一个slog日志实例 level=(DEBUG,INFO,WARN,ERROR); output 日志输出位置
func NewLogger(output io.Writer, opts *slog.HandlerOptions, logType logType) *slog.Logger {
	if output == nil {
		output = os.Stderr
	}
	var handler slog.Handler

	if logType == TextType {
		handler = slog.NewTextHandler(output, opts)
	} else {
		handler = slog.NewJSONHandler(output, opts)
	}

	myHandler := LogHandler{Handler: handler}

	return slog.New(&myHandler)
}

// DefaultOutput 默认日志输出
func DefaultOutput(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    300, // megabytes
		MaxBackups: 30,
		MaxAge:     30,   //days
		Compress:   true, // disabled by default
	}
}

/*
 * 例子
 */

// func demo01() {
// 	opts := &slog.HandlerOptions{
// 		AddSource: true,
// 		Level:     slog.LevelDebug,
// 		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {

// 			// 取相对路径，输出更简短的路径
// 			if a.Key == slog.SourceKey {
// 				// 此时的a.Value 是 slog.Source 指针
// 				ss, ok := a.Value.Any().(*slog.Source)
// 				if !ok || ss.File == "" {
// 					return a
// 				}
// 				var sep = "projectName/"
// 				relativePath := sep + strings.Split(ss.File, sep)[1]
// 				a.Value = slog.StringValue(fmt.Sprintf("%s %d", relativePath, ss.Line))
// 			}

// 			if a.Key == slog.TimeKey {
// 				datetime := a.Value.Time()
// 				a.Value = slog.StringValue(datetime.Format("2006-01-02 15:04:05.000"))
// 			}

// 			return a
// 		},
// 	}
// 	l := NewLogger(DefaultOutput("/logs/access.log"), opts, JsonType)
// 	l.Info("message", "key1", "value1")
// }
