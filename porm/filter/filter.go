package filter

// Filter 过滤器
type Filter interface {
	Type(name string) (interface{}, bool)
	Value(name string, value interface{}) (interface{}, bool)
}

func NewEmptyFieldFilter() Filter {
	return EmptyFieldFilter("")
}

// EmptyFieldFilter 空的字段过滤器
type EmptyFieldFilter string

// Type 类型
func (e EmptyFieldFilter) Type(_ string) (interface{}, bool) {
	return nil, false
}

// Value 值
func (e EmptyFieldFilter) Value(_ string, _ interface{}) (interface{}, bool) {
	return nil, false
}
