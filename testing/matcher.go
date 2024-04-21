package testing

func NewMatcher[A any](name string, match func(a A) bool) Matcher[A] {
	return &matcher[A]{
		name:  name,
		match: match,
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

type matcher[A any] struct {
	name  string
	match func(a A) bool
}

func (m *matcher[A]) Match(actual A) bool {
	return m.match(actual)
}

func (m *matcher[A]) Expected() any {
	return nil
}

func (m *matcher[A]) Negative() bool {
	return false
}

func (m *matcher[A]) Name() string {
	return m.name
}
