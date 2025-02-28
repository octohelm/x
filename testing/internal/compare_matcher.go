package internal

import "fmt"

func NewCompareMatcher[A any, E any](name string, match func(a A, e E) bool) func(e E) Matcher[A] {
	return func(expected E) Matcher[A] {
		return &compareMatcher[A, E]{
			name:     name,
			match:    match,
			expected: expected,
		}
	}
}

type compareMatcher[A any, E any] struct {
	name     string
	match    func(a A, e E) bool
	expected E
}

func (m *compareMatcher[A, E]) Match(actual A) bool {
	return m.match(actual, m.expected)
}

func (m *compareMatcher[A, E]) Negative() bool {
	return false
}

func (m *compareMatcher[A, E]) Name() string {
	return m.name
}

func (m *compareMatcher[A, E]) FormatActual(actual A) string {
	return fmt.Sprintf("%v", actual)
}

var _ ExpectedFormatter = &compareMatcher[any, any]{}

func (m *compareMatcher[A, E]) FormatExpected() string {
	return fmt.Sprintf("%v", m.expected)
}
