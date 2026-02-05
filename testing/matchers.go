package testing

import (
	"reflect"

	reflectx "github.com/octohelm/x/reflect"
	"github.com/octohelm/x/testing/internal/deprecated"
	"github.com/octohelm/x/testing/snapshot"
)

func NotBeNil[T any]() Matcher[T] {
	return Not(BeNil[T]())
}

func BeNil[T any]() Matcher[T] {
	return deprecated.NewMatcher[T]("BeNil", func(a T) bool {
		return any(a) == nil
	})
}

func BeTrue() Matcher[bool] {
	return deprecated.NewMatcher[bool]("BeTrue", func(a bool) bool {
		return a
	})
}

func BeFalse() Matcher[bool] {
	return deprecated.NewMatcher[bool]("BeFalse", func(a bool) bool {
		return !a
	})
}

func NotBeEmpty[T any]() Matcher[T] {
	return Not(BeEmpty[T]())
}

func BeEmpty[T any]() Matcher[T] {
	return deprecated.NewMatcher("BeEmpty", func(a T) bool {
		return reflectx.IsEmptyValue(a)
	})
}

func NotBe[T any](e T) Matcher[T] {
	return Not(Be[T](e))
}

func Be[T any](e T) Matcher[T] {
	return deprecated.NewCompareMatcher[T, T]("Be", func(a T, e T) bool {
		return any(a) == any(e)
	})(e)
}

func NotEqual[T any](e T) Matcher[T] {
	return deprecated.Not(Equal[T](e))
}

func Equal[T any](e T) Matcher[T] {
	return deprecated.NewCompareMatcher[T, T]("Equal", func(a T, e T) bool {
		return reflect.DeepEqual(e, a)
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

func NewSnapshot() *snapshot.Snapshot {
	return &snapshot.Snapshot{}
}

type SnapshotOption = func(m *snapshotMatcher)

func WithWorkDir(wd string) SnapshotOption {
	return func(m *snapshotMatcher) {
		m.wd = wd
	}
}

func MatchSnapshot(name string, options ...SnapshotOption) deprecated.Matcher[*snapshot.Snapshot] {
	c := snapshot.Context{Name: name}
	s, err := c.Load()
	if err != nil {
		panic(err)
	}

	m := &snapshotMatcher{
		expected: s,
	}

	for _, fn := range options {
		fn(m)
	}

	return m
}

type snapshotMatcher struct {
	wd       string
	expected *snapshot.Snapshot
}

func (snapshotMatcher) Action() string {
	return "match snapshot"
}

func (s *snapshotMatcher) Negative() bool {
	return false
}

func (s *snapshotMatcher) Match(a *snapshot.Snapshot) bool {
	return s.expected.Equal(a)
}
