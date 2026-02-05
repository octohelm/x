package deprecated

import (
	"github.com/octohelm/x/testing/lines"
)

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

func (m *compareMatcher[A, E]) Action() string {
	return m.action
}

func (m *compareMatcher[A, E]) Match(actual A) bool {
	return m.match(actual, m.expected)
}

func (m *compareMatcher[A, E]) Negative() bool {
	return false
}

var _ MatcherWithNormalizedExpected = &compareMatcher[string, string]{}

func (m *compareMatcher[A, E]) NormalizedExpected() any {
	switch x := any(m.expected).(type) {
	case lines.Differ:
		return x.Lines()
	default:
		return x
	}
}

var _ MatcherWithActualNormalizer[string] = &compareMatcher[string, string]{}

func (compareMatcher[A, E]) NormalizeActual(a A) any {
	switch x := any(a).(type) {
	case lines.Differ:
		return x.Lines()
	default:
		return x
	}
}
