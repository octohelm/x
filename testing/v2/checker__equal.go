package v2

import (
	"reflect"

	"github.com/octohelm/x/testing/internal"
)

// Equal 返回要求实际值与期望值深度相等的检查器。
func Equal[V any](expect V) ValueChecker[V] {
	return internal.Helper(1, &beChecker[V]{
		be: func(actual V) error {
			if reflect.DeepEqual(expect, actual) {
				return nil
			}
			return &ErrNotEqual{Expect: expect, Got: actual}
		},
	})
}

// NotEqual 返回要求实际值与期望值深度不相等的检查器。
func NotEqual[V any](expect V) ValueChecker[V] {
	return internal.Helper(1, &beChecker[V]{
		be: func(actual V) error {
			if !reflect.DeepEqual(expect, actual) {
				return nil
			}
			return &ErrEqual{NotExpect: expect, Got: actual}
		},
	})
}
