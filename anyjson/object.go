package anyjson

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"strconv"

	"errors"
	"github.com/go-json-experiment/json"
	jsonv1 "github.com/go-json-experiment/json/v1"

	"github.com/go-json-experiment/json/jsontext"
)

type field struct {
	key   string
	value Valuer
}

func (f *field) Set(v Valuer) {
	f.value = v
}

type Object struct {
	props map[string]*node[*field]
	ll    list[*field]
}

func (v *Object) Value() any {
	m := map[string]any{}
	for k, e := range v.props {
		m[k] = e.Value.value.Value()
	}
	return m
}

func (v *Object) Len() int {
	return len(v.props)
}

func (v *Object) KeyValues() iter.Seq2[string, Valuer] {
	return func(yield func(string, Valuer) bool) {
		for el := v.ll.Front; el != nil; el = el.Next {
			if !yield(el.Value.key, el.Value.value) {
				return
			}
		}
	}
}

func (v *Object) Get(key string) (Valuer, bool) {
	if v.props != nil {
		v, ok := v.props[key]
		if ok {
			return v.Value.value, true
		}
	}
	return nil, false
}

func (v *Object) Set(key string, value Valuer) bool {
	if v.props == nil {
		v.props = map[string]*node[*field]{}
	}

	_, alreadyExist := v.props[key]
	if alreadyExist {
		v.props[key].Value.Set(value)
		return false
	}

	element := &field{key: key, value: value}
	v.props[key] = v.ll.PushBack(element)
	return true
}

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

var _ json.UnmarshalerV2 = &Object{}

func (v *Object) UnmarshalJSONV2(d *jsontext.Decoder, options json.Options) error {
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

func (v *Object) UnmarshalJSON(b []byte) error {
	return v.UnmarshalJSONV2(jsontext.NewDecoder(bytes.NewReader(b)), jsonv1.DefaultOptionsV1())
}

func (v *Object) MarshalJSON() ([]byte, error) {
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

func (v *Object) String() string {
	return ToString(v)
}

type list[V any] struct {
	Front, Back *node[V]
}

type node[V any] struct {
	Value      V
	Prev, Next *node[V]
}

func (l *list[V]) PushBack(v V) *node[V] {
	n := &node[V]{
		Value: v,
	}
	l.PushBackNode(n)
	return n
}

func (l *list[V]) PushFront(v V) *node[V] {
	n := &node[V]{
		Value: v,
	}
	l.PushFrontNode(n)
	return n
}

func (l *list[V]) PushBackNode(n *node[V]) {
	n.Next = nil
	n.Prev = l.Back
	if l.Back != nil {
		l.Back.Next = n
	} else {
		l.Front = n
	}
	l.Back = n
}

func (l *list[V]) PushFrontNode(n *node[V]) {
	n.Next = l.Front
	n.Prev = nil
	if l.Front != nil {
		l.Front.Prev = n
	} else {
		l.Back = n
	}
	l.Front = n
}

func (l *list[V]) Remove(n *node[V]) {
	if n.Next != nil {
		n.Next.Prev = n.Prev
	} else {
		l.Back = n.Prev
	}
	if n.Prev != nil {
		n.Prev.Next = n.Next
	} else {
		l.Front = n.Next
	}
}
