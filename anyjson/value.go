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
	"github.com/octohelm/x/ptr"
	"golang.org/x/sync/errgroup"
)

func Equal(a Valuer, b Valuer) bool {
	return reflect.DeepEqual(a.Value(), b.Value())
}

func MustFromValue(value any) Valuer {
	x, err := FromValue(value)
	if err != nil {
		panic(err)
	}
	return x
}

func As[T Valuer](valuer T, target any) error {
	raw, err := valuer.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, target)
}

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
		return json.MarshalWrite(w, value, jsonv1.OmitEmptyWithLegacyDefinition(true))
	})

	p := &payload{}

	eg.Go(func() error {
		defer r.Close()
		return json.UnmarshalRead(r, p, jsonv1.OmitEmptyWithLegacyDefinition(true))
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return p.Valuer, nil
}

type payload struct {
	Valuer
}

func (p *payload) UnmarshalJSONFrom(decoder *jsontext.Decoder, options json.Options) error {
	v, err := FromJSONTextDecoder(decoder)
	if err != nil {
		return err
	}
	p.Valuer = v
	return nil
}

type Valuer interface {
	json.Marshaler

	fmt.Stringer

	Value() any
}

func ToString(valuer Valuer) string {
	data, _ := valuer.MarshalJSON()
	return string(data)
}

func FromJSONTextDecoder(decoder *jsontext.Decoder) (Valuer, error) {
	switch decoder.PeekKind() {
	case '{':
		o := &Object{}
		if err := o.UnmarshalJSONFrom(decoder, json.DefaultOptionsV2()); err != nil {
			return nil, err
		}
		return o, nil
	case '[':
		arr := &Array{}
		if err := arr.UnmarshalJSONFrom(decoder, json.DefaultOptionsV2()); err != nil {
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
		return &Boolean{raw: value.Clone(), value: ptr.Ptr(false)}, nil
	case 't':
		value, err := decoder.ReadValue()
		if err != nil {
			return nil, err
		}
		return &Boolean{raw: value.Clone(), value: ptr.Ptr(true)}, nil
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
