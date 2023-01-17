package forms

// 根据字段存储错误信息，一个字段可以有多个错误
type errMaps map[string][]string

// Add 根据字段往 errMaps 插入一条错误
func (em errMaps) Add(field, message string) {
	em[field] = append(em[field], message)
}

// Get 根据字段从 errMaps 获取一个错误，如果没有返回空串
func (em errMaps) Get(field string) string {
	errs := em[field]
	if len(errs) == 0 {
		return ""
	}
	return errs[0]
}
