package anyjson

import (
	"bytes"
	"strconv"
)

// BooleanOf 创建一个布尔值节点。
func BooleanOf(b bool) *Boolean {
	return &Boolean{value: &b}
}

// Boolean 表示 JSON 布尔值。
type Boolean struct {
	raw   []byte
	value *bool
}

func (v *Boolean) MarshalJSON() ([]byte, error) {
	if v.raw == nil && v.value != nil {
		v.raw = []byte(strconv.FormatBool(*v.value))
	}
	return v.raw, nil
}

func (v *Boolean) Value() any {
	if v.value == nil {
		v.value = new(bytes.Equal(v.raw, []byte("true")))
	}
	return *v.value
}

// String 返回布尔值的 JSON 文本表示。
func (v *Boolean) String() string {
	return ToString(v)
}
