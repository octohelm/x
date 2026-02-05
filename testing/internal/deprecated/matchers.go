package deprecated

import (
	"reflect"

	reflectx "github.com/octohelm/x/reflect"
)

type Equatable[T any] interface {
	Equal(T) bool
}

func Equal[T any](e T) Matcher[T] {
	return NewCompareMatcher[T, T]("equal", func(a T, e T) bool {
		if x, ok := any(e).(Equatable[T]); ok {
			return x.Equal(a)
		}
		return reflect.DeepEqual(a, e)
	})(e)
}

func NotEqual[T any](e T) Matcher[T] {
	return Not(Equal[T](e))
}

func BeTrue() Matcher[bool] {
	return NewMatcher[bool]("be true", func(a bool) bool {
		return a
	})
}

func BeFalse() Matcher[bool] {
	return NewMatcher[bool]("be false", func(a bool) bool {
		return !a
	})
}

func BeNil[T any]() Matcher[T] {
	return NewMatcher[T]("be nil", func(a T) bool {
		return any(a) == nil
	})
}

func BeZero[T any]() Matcher[T] {
	return NewMatcher[T]("be zero", func(a T) bool {
		return reflectx.IsZero(a)
	})
}
