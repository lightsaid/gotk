package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// 正则表达式部分参考：
// github.com/asaskevich/govalidator
var (
	EmailRX = regexp.MustCompile("^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")
	PhonePX = regexp.MustCompile(`^1[3-9]\d{9}$`)
)

// Form 自定义Form结构，包含 url.Values 和 错误信息errMaps
type Form struct {
	url.Values
	Errors errMaps
}

// New 创建一个Form实例
func New(data url.Values) *Form {
	return &Form{
		data,
		errMaps(map[string][]string{}),
	}
}

// Valid 判断验证是否通过，true 通过
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Required 必填字段
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		val := f.Get(field)
		if strings.TrimSpace(val) == "" {
			f.Errors.Add(field, field+" 不能为空")
		}
	}
}

// RequiredForMsg 如果 field 字段为空，则添加 msg 错误消息
func (f *Form) RequiredForMsg(field, msg string) {
	if strings.TrimSpace(f.Get(field)) == "" {
		f.Errors.Add(field, msg)
	}
}

// MaxLength 字符最大长度不能超过max
func (f *Form) MaxLength(field string, max int, msgs ...string) {
	val := f.Get(field)
	if val == "" {
		return
	}

	if utf8.RuneCountInString(val) > max {
		errMsg := fmt.Sprintf("%s 长度不能超过 %d", field, max)
		if len(msgs) > 0 {
			errMsg = msgs[0]
		}
		f.Errors.Add(field, errMsg)
	}
}

// MinLength 字符最小长度不能小于min
func (f *Form) MinLength(field string, min int, msgs ...string) {
	if utf8.RuneCountInString(f.Get(field)) < min {
		errMsg := fmt.Sprintf("%s 长度不能小于 %d", field, min)
		if len(msgs) > 0 {
			errMsg = msgs[0]
		}
		f.Errors.Add(field, errMsg)
	}
}

// Check expr 表达式如果不成立，则将field 和 msg 添加到错误信息 Errors
func (f *Form) Check(expr bool, field, msg string) {
	if !expr {
		f.Errors.Add(field, msg)
	}
}

// IsEmail 是否是email
func (f *Form) IsEmail(field, msg string) {
	if !f.Matches(f.Get(field), EmailRX) {
		f.Errors.Add(field, msg)
	}
}

// IsPhone 是否是手机
func (f *Form) IsPhone(field, msg string) {
	if !f.Matches(f.Get(field), PhonePX) {
		f.Errors.Add(field, msg)
	}
}

// Matches value是否匹配正则表达式 rgx
func (f *Form) Matches(value string, rgx *regexp.Regexp) bool {
	return rgx.MatchString(value)
}
