package anyjson

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"iter"
	"strconv"
	"sync/atomic"

	jsonv1 "github.com/go-json-experiment/json/v1"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/octohelm/x/container/list"
)

type field struct {
	key   string
	value Valuer
}

func (f *field) Set(v Valuer) {
	f.value = v
}

// Object 表示保持键插入顺序的 JSON 对象。
type Object struct {
	props     map[string]*list.Element[*field]
	ll        list.List[*field]
	initialed uint32
}

func (v *Object) init() {
	if atomic.SwapUint32(&v.initialed, 1) == 0 {
		v.props = map[string]*list.Element[*field]{}
		v.ll.Init()
	}
}

func (v *Object) Value() any {
	m := map[string]any{}
	for k, e := range v.props {
		m[k] = e.Value.value.Value()
	}
	return m
}

// Len 返回对象中的键值对数量。
func (v *Object) Len() int {
	return len(v.props)
}

// KeyValues 以插入顺序遍历对象中的键和值。
func (v *Object) KeyValues() iter.Seq2[string, Valuer] {
	return func(yield func(string, Valuer) bool) {
		for el := v.ll.Front(); el != nil; el = el.Next() {
			if !yield(el.Value.key, el.Value.value) {
				return
			}
		}
	}
}

// Get 返回指定键对应的值。
func (v *Object) Get(key string) (Valuer, bool) {
	if v.props != nil {
		v, ok := v.props[key]
		if ok {
			return v.Value.value, true
		}
	}
	return nil, false
}

// Set 设置键对应的值。
//
// 当 key 首次写入时返回 true；更新已有键时返回 false。
func (v *Object) Set(key string, value Valuer) bool {
	v.init()

	_, alreadyExist := v.props[key]
	if alreadyExist {
		v.props[key].Value.Set(value)
		return false
	}

	element := &field{key: key, value: value}
	v.props[key] = v.ll.PushBack(element)
	return true
}

// Delete 删除指定键，并返回是否实际删除了元素。
func (v *Object) Delete(key string) (didDelete bool) {
	if v.props == nil {
		return false
	}

	element, ok := v.props[key]
	if ok {
		v.ll.Remove(element)

		delete(v.props, key)
	}
	return ok
}

var _ json.UnmarshalerFrom = &Object{}

func (v *Object) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	t, err := d.ReadToken()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	kind := t.Kind()

	if kind != '{' {
		return &json.SemanticError{
			JSONPointer: d.StackPointer(),
			Err:         fmt.Errorf("object should starts with `{`, but got `%s`", kind),
		}
	}

	if v == nil {
		*v = Object{}
	}

	for kind := d.PeekKind(); kind != '}'; kind = d.PeekKind() {
		k, err := d.ReadValue()
		if err != nil {
			return err
		}

		key, err := strconv.Unquote(string(k))
		if err != nil {
			return &json.SemanticError{
				JSONPointer: d.StackPointer(),
				Err:         errors.New("key should be quoted string"),
			}
		}

		value, err := FromJSONTextDecoder(d)
		if err != nil {
			return err
		}

		v.Set(key, value)
	}

	// read the close '}'
	if _, err := d.ReadToken(); err != nil {
		if err != io.EOF {
			return nil
		}
		return err
	}
	return nil
}

// UnmarshalJSON 将 JSON 文本解码为对象。
func (v *Object) UnmarshalJSON(b []byte) error {
	o := &Object{}

	if err := o.UnmarshalJSONFrom(jsontext.NewDecoder(bytes.NewReader(b), jsonv1.OmitEmptyWithLegacySemantics(true))); err != nil {
		return err
	}

	*v = *o

	return nil
}

// MarshalJSON 将对象编码为 JSON 文本，并保留键的插入顺序。
func (v Object) MarshalJSON() ([]byte, error) {
	b := bytes.NewBuffer(nil)

	b.WriteString("{")

	idx := 0
	for k, v := range v.KeyValues() {
		if idx > 0 {
			b.WriteString(",")
		}

		b.WriteString(strconv.Quote(k))
		b.WriteString(":")
		raw, err := v.MarshalJSON()
		if err != nil {
			return []byte{}, err
		}
		b.Write(raw)
		idx++
	}

	b.WriteString("}")

	return b.Bytes(), nil
}

// String 返回对象的 JSON 文本表示。
func (v *Object) String() string {
	return ToString(v)
}
