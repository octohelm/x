package testing

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

func (m *compareMatcher[A, E]) Expected() any {
	return m.expected
}
