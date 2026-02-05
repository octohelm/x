package v2

import (
	"reflect"

	"github.com/octohelm/x/testing/internal"
)

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
