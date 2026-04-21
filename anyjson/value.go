package anyjson

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	jsonv1 "github.com/go-json-experiment/json/v1"
	"golang.org/x/sync/errgroup"
)

// Equal 判断两个 Valuer 在原生 Go 值层面是否相等。
func Equal(a Valuer, b Valuer) bool {
	return reflect.DeepEqual(a.Value(), b.Value())
}

// MustFromValue 将任意 Go 值转换为 Valuer，失败时 panic。
func MustFromValue(value any) Valuer {
	x, err := FromValue(value)
	if err != nil {
		panic(err)
	}
	return x
}

// As 将 Valuer 反序列化到 target 指向的目标值中。
func As[T Valuer](valuer T, target any) error {
	raw, err := valuer.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, target)
}

// FromValue 将任意 Go 值转换为 anyjson 的值表示。
//
// 它会保留对象、数组、标量和 nil 的结构语义，便于后续做 diff、merge 或 transform。
func FromValue(value any) (Valuer, error) {
	if value == nil {
		return &Null{}, nil
	}

	switch x := value.(type) {
	case []any:
		arr := &Array{}
		for _, e := range x {
			item, err := FromValue(e)
			if err != nil {
				return nil, err
			}
			arr.Append(item)
		}
		return arr, nil
	case map[string]any:
		o := &Object{}
		keys := make([]string, 0, len(x))
		for key := range x {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			propValue, err := FromValue(x[key])
			if err != nil {
				return nil, err
			}
			o.Set(key, propValue)
		}
		return o, nil
	case string:
		return &String{value: &x}, nil
	case bool:
		return &Boolean{value: &x}, nil
	case int:
		return NumberOf(x), nil
	case int8:
		return NumberOf(x), nil
	case int16:
		return NumberOf(x), nil
	case int32:
		return NumberOf(x), nil
	case int64:
		return NumberOf(x), nil
	case uint:
		return NumberOf(x), nil
	case uint8:
		return NumberOf(x), nil
	case uint16:
		return NumberOf(x), nil
	case uint32:
		return NumberOf(x), nil
	case uint64:
		return NumberOf(x), nil
	case float32:
		return NumberOf(x), nil
	case float64:
		return NumberOf(x), nil
	}

	r, w := io.Pipe()

	eg := &errgroup.Group{}

	eg.Go(func() error {
		defer w.Close()
		return json.MarshalWrite(w, value, jsonv1.OmitEmptyWithLegacySemantics(true))
	})

	p := &payload{}

	eg.Go(func() error {
		defer r.Close()
		return json.UnmarshalRead(r, p, jsonv1.OmitEmptyWithLegacySemantics(true))
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return p.Valuer, nil
}

type payload struct {
	Valuer
}

var _ json.UnmarshalerFrom = &payload{}

func (p *payload) UnmarshalJSONFrom(decoder *jsontext.Decoder) error {
	v, err := FromJSONTextDecoder(decoder)
	if err != nil {
		return err
	}
	p.Valuer = v
	return nil
}

// Valuer 表示 anyjson 中统一的 JSON 值抽象。
//
// 实现类型需要同时支持 JSON 编码、字符串表示和原生 Go 值读取。
type Valuer interface {
	json.Marshaler

	fmt.Stringer

	Value() any
}

// ToString 以 JSON 文本形式返回 valuer 的字符串表示。
func ToString(valuer Valuer) string {
	data, _ := valuer.MarshalJSON()
	return string(data)
}

// FromJSONTextDecoder 从 jsontext.Decoder 当前位置读取一个 JSON 值并转换为 Valuer。
func FromJSONTextDecoder(decoder *jsontext.Decoder) (Valuer, error) {
	switch decoder.PeekKind() {
	case '{':
		o := &Object{}
		if err := o.UnmarshalJSONFrom(decoder); err != nil {
			return nil, err
		}
		return o, nil
	case '[':
		arr := &Array{}
		if err := arr.UnmarshalJSONFrom(decoder); err != nil {
			return nil, err
		}
		return arr, nil
	case 'n':
		_, err := decoder.ReadValue()
		if err != nil {
			return nil, err
		}
		return &Null{}, nil
	case 'f':
		value, err := decoder.ReadValue()
		if err != nil {
			return nil, err
		}
		return &Boolean{raw: value.Clone(), value: new(false)}, nil
	case 't':
		value, err := decoder.ReadValue()
		if err != nil {
			return nil, err
		}
		return &Boolean{raw: value.Clone(), value: new(true)}, nil
	case '"':
		value, err := decoder.ReadValue()
		if err != nil {
			return nil, err
		}
		return &String{raw: value.Clone()}, nil
	case '0':
		value, err := decoder.ReadValue()
		if err != nil {
			return nil, err
		}
		if bytes.Contains(value, []byte(".")) {
			return &Number[float64]{raw: value.Clone()}, nil
		}
		return &Number[int]{raw: value.Clone()}, nil
	}
	return nil, nil
}
