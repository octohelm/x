package internal

type Matcher[A any] interface {
	Action() string
	Match(actual A) bool
	Negative() bool
}

type MatcherWithActualNormalizer[T any] interface {
	NormalizeActual(a T) any
}

type MatcherWithNormalizedExpected interface {
	NormalizedExpected() any
}

func NewMatcher[A any](action string, match func(a A) bool) Matcher[A] {
	return &matcher[A]{
		action: action,
		match:  match,
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

type matcher[A any] struct {
	action string
	match  func(a A) bool
}

func (m *matcher[A]) Match(actual A) bool {
	return m.match(actual)
}

func (m *matcher[A]) Negative() bool {
	return false
}

func (m *matcher[A]) Action() string {
	return m.action
}
