package anyjson

import (
	"bytes"
	"errors"
	"io"
	"iter"

	"github.com/go-json-experiment/json/jsontext"
	jsonv1 "github.com/go-json-experiment/json/v1"
)

type Array struct {
	items []Valuer
}

func (v *Array) Value() any {
	list := make([]any, len(v.items))
	for i := range list {
		list[i] = v.items[i].Value()
	}
	return list
}

func (v *Array) Len() int {
	return len(v.items)
}

func (v *Array) Values() iter.Seq[Valuer] {
	return func(yield func(Valuer) bool) {
		for _, v := range v.items {
			if !yield(v) {
				return
			}
		}
	}
}

func (v *Array) IndexedValues() iter.Seq2[int, Valuer] {
	return func(yield func(int, Valuer) bool) {
		for i, v := range v.items {
			if !yield(i, v) {
				return
			}
		}
	}
}

func (v *Array) UnmarshalJSONFrom(d *jsontext.Decoder, v1 jsonv1.Options) error {
	if v == nil {
		*v = Array{}
	}

	t, err := d.ReadToken()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	if t.Kind() != '[' {
		return errors.New("v should starts with `[`")
	}

	for kind := d.PeekKind(); kind != ']'; kind = d.PeekKind() {
		value, err := FromJSONTextDecoder(d)
		if err != nil {
			return err
		}
		v.items = append(v.items, value)
	}

	// read the close ']'
	if _, err := d.ReadToken(); err != nil {
		if err != io.EOF {
			return nil
		}
		return err
	}
	return nil
}

func (v *Array) UnmarshalJSON(b []byte) error {
	d := jsontext.NewDecoder(bytes.NewReader(b))
	return v.UnmarshalJSONFrom(d, jsonv1.DefaultOptionsV1())
}

func (v *Array) MarshalJSON() ([]byte, error) {
	b := bytes.NewBuffer(nil)

	b.WriteString("[")

	for idx, v := range v.items {
		if idx > 0 {
			b.WriteString(",")
		}

		raw, err := v.MarshalJSON()
		if err != nil {
			return []byte{}, err
		}
		b.Write(raw)
		idx++
	}

	b.WriteString("]")

	return b.Bytes(), nil
}

func (v *Array) String() string {
	return ToString(v)
}

func (v *Array) Append(item Valuer) {
	v.items = append(v.items, item)
}

func (v *Array) Index(i int) (Valuer, bool) {
	if i < 0 || i >= len(v.items) {
		return nil, false
	}

	return v.items[i], true
}
