package testing

import "reflect"

func Be[A any](e A) Matcher[A] {
	return MatcherWith[A, A](func(a A, e A) bool {
		return any(a) == any(e)
	}, "Be")(e)
}

func Equal[A any](e A) Matcher[A] {
	return MatcherWith[A, A](func(a A, e A) bool {
		return reflect.DeepEqual(a, e)
	}, "Equal")(e)
}

func HaveCap[A any](c int) Matcher[A] {
	return MatcherWith[A, int](func(a A, c int) bool {
		return reflect.ValueOf(a).Cap() == c
	}, "HaveCap")(c)
}

func HaveLen[A any](c int) Matcher[A] {
	return MatcherWith[A, int](func(a A, c int) bool {
		return reflect.ValueOf(a).Len() == c
	}, "HaveLen")(c)
}

func MatcherWith[A any, E any](match func(a A, e E) bool, name string) func(e E) Matcher[A] {
	return func(expected E) Matcher[A] {
		return &matcher[A, E]{
			name:     name,
			match:    match,
			expected: expected,
		}
	}
}

func Not[A any](matcher Matcher[A]) Matcher[A] {
	return &negativeMatcher[A]{
		Matcher: matcher,
	}
}

type negativeMatcher[A any] struct {
	Matcher[A]
}

func (m *negativeMatcher[A]) Negative() bool {
	return true
}

type Matcher[A any] interface {
	Name() string
	Expected() any
	Negative() bool
	Match(actual A) bool
}

type matcher[A any, E any] struct {
	name     string
	match    func(a A, e E) bool
	expected E
}

func (m *matcher[A, E]) Match(actual A) bool {
	return m.match(actual, m.expected)
}

func (m *matcher[A, E]) Negative() bool {
	return false
}

func (m *matcher[A, E]) Name() string {
	return m.name
}

func (m *matcher[A, E]) Expected() any {
	return m.expected
}
