package gotk

import (
	"fmt"
	"regexp"
)

// Validator 验证器，Error 保存字段和对应的错误信息
type Validator struct {
	Errors map[string]string
}

// New 创建一个验证器实例
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Valid 验证是否通过
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError 添加一个错误信息
func (v *Validator) AddError(field, message string) {
	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = message
	}
}

// MinLength 最小长度, 如果不满足条件，添加错误，返回false
func (v *Validator) MinLength(field string, length int) bool {
	x := v.Errors[field]
	if len([]rune(x)) < length {
		v.AddError(field, fmt.Sprintf("%s 字段长度必须 >= %d", field, length))
		return false
	}
	return true
}

// MaxLength 最大长度, 如果不满足条件，添加错误，返回false
func (v *Validator) MaxLength(field string, length int) bool {
	x := v.Errors[field]
	if len([]rune(x)) >= length {
		v.AddError(field, fmt.Sprintf("%s 字段长度必须 <= %d", field, length))
		return false
	}
	return true
}

// Check检查，当expr表达式为false添加错误，
// 因此expr可以理解为满足条件的表达式
func (v *Validator) Check(expr bool, field, message string) {
	if !expr {
		v.AddError(field, message)
	}
}

// Matches value 是否满足rx正则
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

//===============================================================

// 辅助验证

// Verifiyer定义一个Verifiy校验接口方法，
// 配合DoVerifiy使用可以省去validator.New(),Valid()步骤
// Verifiy 执行验证字段方法，将错误添加到 *Validator 里
type Verifiyer interface {
	Verifiy(*Validator)
}

// DoVerifiy 执行实现Verifiyer接口对象，并返回*Validator
func DoVerifiy(obj Verifiyer) (*Validator, bool) {
	v := New()
	obj.Verifiy(v)
	return v, v.Valid()
}
