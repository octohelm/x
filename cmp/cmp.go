package cmp

import (
	"cmp"
	"errors"
	"fmt"
	"iter"
	"reflect"

	"github.com/go-json-experiment/json/jsontext"
)

// True 返回要求实际值为 true 的谓词。
func True() func(a bool) error {
	return Eq(true)
}

// False 返回要求实际值为 false 的谓词。
func False() func(a bool) error {
	return Eq(false)
}

// Eq 返回要求实际值等于期望值的谓词。
func Eq[V comparable](e V) func(a V) error {
	return func(a V) error {
		if a == e {
			return nil
		}
		return &ErrCondition{Op: "==", Expect: e, Actual: a}
	}
}

// Neq 返回要求实际值不等于期望值的谓词。
func Neq[V comparable](e V) func(a V) error {
	return func(a V) error {
		if a != e {
			return nil
		}
		return &ErrCondition{Op: "!=", Expect: e, Actual: a}
	}
}

// Gt 返回要求实际值大于期望值的谓词。
func Gt[V cmp.Ordered](e V) func(a V) error {
	return func(a V) error {
		if a > e {
			return nil
		}
		return &ErrCondition{Op: ">", Expect: e, Actual: a}
	}
}

// Gte 返回要求实际值大于等于期望值的谓词。
func Gte[V cmp.Ordered](e V) func(a V) error {
	return func(a V) error {
		if a >= e {
			return nil
		}
		return &ErrCondition{Op: ">=", Expect: e, Actual: a}
	}
}

// Lt 返回要求实际值小于期望值的谓词。
func Lt[V cmp.Ordered](e V) func(a V) error {
	return func(a V) error {
		if a < e {
			return nil
		}
		return &ErrCondition{Op: "<", Expect: e, Actual: a}
	}
}

// Lte 返回要求实际值小于等于期望值的谓词。
func Lte[V cmp.Ordered](e V) func(a V) error {
	return func(a V) error {
		if a <= e {
			return nil
		}
		return &ErrCondition{Op: "<=", Expect: e, Actual: a}
	}
}

// Nil 返回要求实际值为 nil 的谓词。
func Nil[V any]() func(a V) error {
	return func(a V) error {
		if isNil(a) {
			return nil
		}
		return &ErrState{State: "nil", Actual: a}
	}
}

// NotNil 返回要求实际值非 nil 的谓词。
func NotNil[V any]() func(a V) error {
	return func(a V) error {
		if !isNil(a) {
			return nil
		}
		return &ErrState{State: "not nil", Actual: a}
	}
}

func isNil(a any) bool {
	rv := reflect.ValueOf(a)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		if rv.IsNil() {
			return true
		}
	case reflect.Invalid:
		return true
	default:
	}
	return false
}

// Zero 返回要求实际值为零值的谓词。
func Zero[V any]() func(a V) error {
	return func(a V) error {
		if reflect.ValueOf(a).IsZero() {
			return nil
		}
		return &ErrState{State: "zero", Actual: a}
	}
}

// NotZero 返回要求实际值不为零值的谓词。
func NotZero[V any]() func(a V) error {
	return func(a V) error {
		if !reflect.ValueOf(a).IsZero() {
			return nil
		}
		return &ErrState{State: "not zero", Actual: a}
	}
}

// Len 返回针对长度值的谓词。
//
// e 可以是固定长度，也可以是进一步校验长度的谓词。
func Len[V any, E int | func(int) error](e E) func(a V) error {
	return func(a V) error {
		t := reflect.TypeOf(a)

		var n int

		switch t.Kind() {
		case reflect.Slice, reflect.Map, reflect.Chan, reflect.Array, reflect.String:
			n = reflect.ValueOf(a).Len()
		default:
			return &ErrState{State: "lengthable", Actual: a}
		}

		var err error
		switch x := any(e).(type) {
		case int:
			err = Eq(x)(n)
		case func(int) error:
			err = x(n)
		}

		if err != nil {
			return &ErrCheck{
				Topic:  "len",
				Err:    err,
				Actual: a,
			}
		}

		return nil
	}
}

// Every 返回要求序列中每个元素都满足谓词的谓词。
func Every[V any](p func(V) error) func(seq iter.Seq[V]) error {
	return func(seq iter.Seq[V]) error {
		i := 0
		for item := range seq {
			if err := p(item); err != nil {
				return wrap(err, "elem", fmt.Sprintf("%d", i), item)
			}
			i++
		}
		return nil
	}
}

// Some 返回要求序列中至少一个元素满足谓词的谓词。
func Some[V any](p func(V) error) func(seq iter.Seq[V]) error {
	return func(seq iter.Seq[V]) error {
		var lastErr error
		for item := range seq {
			if err := p(item); err == nil {
				return nil
			} else {
				lastErr = err
			}
		}
		return &ErrCheck{
			Topic: "elem",
			Err:   fmt.Errorf("none of the elements satisfy the predicate (error: %w)", lastErr),
		}
	}
}

func wrap(err error, topic string, tok string, actual any) error {
	if err == nil {
		return nil
	}

	next := &ErrCheck{
		Topic:  topic,
		Err:    err,
		Actual: actual,
	}

	if child, ok := errors.AsType[*ErrCheck](err); ok {
		next.Pointer = jsontext.Pointer("").AppendToken(tok) + child.Pointer
		next.Err = child.Err
		next.Topic = child.Topic
	} else {
		next.Pointer = jsontext.Pointer("").AppendToken(tok)
	}

	return next
}
