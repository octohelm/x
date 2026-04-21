package anyjson

import (
	"strconv"

	"github.com/go-json-experiment/json"
)

// StringOf 创建一个字符串值节点。
func StringOf(v string) *String {
	return &String{value: &v}
}

// String 表示 JSON 字符串值。
type String struct {
	value *string
	raw   []byte
}

func (v *String) MarshalJSON() ([]byte, error) {
	if v.raw == nil && v.value != nil {
		v.raw, _ = json.Marshal(v.value)
	}
	return v.raw, nil
}

func (v *String) Value() any {
	if v.value == nil {
		s, _ := strconv.Unquote(string(v.raw))
		v.value = &s
	}
	return *v.value
}

// String 返回字符串值的 JSON 文本表示。
func (v *String) String() string {
	return ToString(v)
}
