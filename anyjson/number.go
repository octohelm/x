package anyjson

import (
	"strconv"

	"github.com/go-json-experiment/json"
	"github.com/octohelm/x/ptr"
)

func NumberOf[T number](n T) *Number[T] {
	return &Number[T]{
		value: &n,
	}
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type Number[T number] struct {
	value *T
	raw   []byte
}

func (v *Number[T]) MarshalJSON() ([]byte, error) {
	if v.raw == nil && v.value != nil {
		v.raw, _ = json.Marshal(v.value)
	}
	return v.raw, nil
}

func (v *Number[T]) Value() any {
	if v.value == nil {
		i, err := strconv.ParseInt(string(v.raw), 10, 64)
		if err == nil {
			v.value = ptr.Ptr(T(i))
		} else {
			f, err := strconv.ParseFloat(string(v.raw), 64)
			if err == nil {
				v.value = ptr.Ptr(T(f))
			}
		}
	}
	if v := v.value; v != nil {
		return *v
	}
	return nil
}

func (v *Number[T]) String() string {
	return ToString(v)
}
