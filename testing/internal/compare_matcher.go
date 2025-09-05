package internal

func NewCompareMatcher[A any, E any](action string, match func(a A, e E) bool) func(e E) Matcher[A] {
	return func(expected E) Matcher[A] {
		return &compareMatcher[A, E]{
			action:   action,
			match:    match,
			expected: expected,
		}
	}
}

type compareMatcher[A any, E any] struct {
	action   string
	match    func(a A, e E) bool
	expected E
}

func (m *compareMatcher[A, E]) Match(actual A) bool {
	return m.match(actual, m.expected)
}

func (m *compareMatcher[A, E]) Negative() bool {
	return false
}

func (m *compareMatcher[A, E]) Action() string {
	return m.action
}

var _ MatcherWithNormalizedExpected = &compareMatcher[string, string]{}

func (m *compareMatcher[A, E]) NormalizedExpected() any {
	return m.expected
}
