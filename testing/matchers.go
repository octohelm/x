package testing

import (
	"reflect"
)

func Be[A any](a A, e A) bool {
	return any(a) == any(e)
}

func Equal[A any](a A, e A) bool {
	return reflect.DeepEqual(a, e)
}

func HaveCap[A any](a A, c int) bool {
	return reflect.ValueOf(a).Cap() == c
}

func HaveLen[A any](a A, c int) bool {
	return reflect.ValueOf(a).Len() == c
}
