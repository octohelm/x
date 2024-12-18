package anyjson

import (
	"strconv"

	"github.com/go-json-experiment/json"
)

func StringOf(v string) *String {
	return &String{value: &v}
}

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

func (v *String) String() string {
	return ToString(v)
}
