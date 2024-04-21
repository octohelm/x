package testing

import (
	"reflect"

	reflectx "github.com/octohelm/x/reflect"
)

func NotBeNil[T any]() Matcher[T] {
	return Not(BeNil[T]())
}

func BeNil[T any]() Matcher[T] {
	return NewMatcher[T]("BeNil", func(a T) bool {
		return any(a) == nil
	})
}

func BeTrue() Matcher[bool] {
	return NewMatcher[bool]("BeTrue", func(a bool) bool {
		return a
	})
}

func BeFalse() Matcher[bool] {
	return NewMatcher[bool]("BeFalse", func(a bool) bool {
		return !a
	})
}

func NotBeEmpty[T any]() Matcher[T] {
	return Not(BeEmpty[T]())
}

func BeEmpty[T any]() Matcher[T] {
	return NewMatcher("BeEmpty", func(a T) bool {
		return reflectx.IsEmptyValue(a)
	})
}

func NotBe[T any](e T) Matcher[T] {
	return Not(Be[T](e))
}

func Be[T any](e T) Matcher[T] {
	return NewCompareMatcher[T, T]("Be", func(a T, e T) bool {
		return any(a) == any(e)
	})(e)
}

func NotEqual[T any](e T) Matcher[T] {
	return Not(Equal[T](e))
}

func Equal[T any](e T) Matcher[T] {
	return NewCompareMatcher[T, T]("Equal", func(a T, e T) bool {
		return reflect.DeepEqual(a, e)
	})(e)
}

func NotHaveCap[T any](c int) Matcher[T] {
	return Not(HaveCap[T](c))
}

func HaveCap[T any](c int) Matcher[T] {
	return NewCompareMatcher[T, int]("HaveCap", func(a T, c int) bool {
		return reflect.ValueOf(a).Cap() == c
	})(c)
}

func NotHaveLen[T any](c int) Matcher[T] {
	return Not(HaveLen[T](c))
}

func HaveLen[T any](c int) Matcher[T] {
	return NewCompareMatcher[T, int]("HaveLen", func(a T, c int) bool {
		return reflect.ValueOf(a).Len() == c
	})(c)
}
