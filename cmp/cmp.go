package cmp

import (
	"cmp"
	"reflect"
)

func True() func(a bool) error {
	return Eq(true)
}

func False() func(a bool) error {
	return Eq(false)
}

func Eq[V comparable](e V) func(a V) error {
	return func(a V) error {
		if a == e {
			return nil
		}
		return &ErrCondition{Op: "==", Expect: e, Actual: a}
	}
}

func Neq[V comparable](e V) func(a V) error {
	return func(a V) error {
		if a != e {
			return nil
		}
		return &ErrCondition{Op: "!=", Expect: e, Actual: a}
	}
}

func Gt[V cmp.Ordered](e V) func(a V) error {
	return func(a V) error {
		if a > e {
			return nil
		}
		return &ErrCondition{Op: ">", Expect: e, Actual: a}
	}
}

func Gte[V cmp.Ordered](e V) func(a V) error {
	return func(a V) error {
		if a >= e {
			return nil
		}
		return &ErrCondition{Op: ">=", Expect: e, Actual: a}
	}
}

func Lt[V cmp.Ordered](e V) func(a V) error {
	return func(a V) error {
		if a < e {
			return nil
		}
		return &ErrCondition{Op: "<", Expect: e, Actual: a}
	}
}

func Lte[V cmp.Ordered](e V) func(a V) error {
	return func(a V) error {
		if a <= e {
			return nil
		}
		return &ErrCondition{Op: "<=", Expect: e, Actual: a}
	}
}

func Nil[V any]() func(a V) error {
	return func(a V) error {
		if isNil(a) {
			return nil
		}
		return &ErrState{State: "nil", Actual: a}
	}
}

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

func Zero[V any]() func(a V) error {
	return func(a V) error {
		if reflect.ValueOf(a).IsZero() {
			return nil
		}
		return &ErrState{State: "zero", Actual: a}
	}
}

func NotZero[V any]() func(a V) error {
	return func(a V) error {
		if !reflect.ValueOf(a).IsZero() {
			return nil
		}
		return &ErrState{State: "not zero", Actual: a}
	}
}

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
