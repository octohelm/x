package bdd

import (
	"errors"
	"fmt"
	"iter"
	"slices"
	"testing"

	"github.com/octohelm/x/testing/internal"
	"github.com/octohelm/x/testing/snapshot"
)

func SliceHaveLen[Slice ~[]E, E any](expect int, actual Slice) Checker {
	matcher := internal.NewCompareMatcher(fmt.Sprintf("have length"), func(a Slice, n int) bool {
		return len(a) == n
	})(expect)
	return asChecker(matcher, actual)
}

func ErrorAs[V error](expect *V, err error) Checker {
	matcher := internal.NewMatcher[error](fmt.Sprintf("as %T", *new(T)), func(a error) bool {
		if expect == nil {
			return false
		}
		return errors.As(a, expect)
	})

	return asChecker(matcher, err)
}

func ErrorIs(expect error, err error) Checker {
	matcher := internal.NewMatcher[error](fmt.Sprintf("is %v", expect), func(a error) bool {
		return errors.Is(a, expect)
	})
	return asChecker(matcher, err)
}

func NoError(err error) Checker {
	matcher := internal.NewMatcher[error]("no error", func(e error) bool {
		return e == nil
	})
	return asChecker(matcher, err)
}

func Zero[V any](actual V) Checker {
	matcher := internal.BeZero[V]()
	return asChecker(matcher, actual)
}

func Nil[V any](actual V) Checker {
	matcher := internal.BeNil[V]()
	return asChecker(matcher, actual)
}

func True(actual bool) Checker {
	matcher := internal.BeTrue()
	return asChecker(matcher, actual)
}

func False(actual bool) Checker {
	matcher := internal.BeFalse()
	return asChecker(matcher, actual)
}

func Equal[V any](expect V, actual V) Checker {
	matcher := internal.Equal(expect)
	return asChecker(matcher, actual)
}

func EqualSeq[V any](expect iter.Seq[V], actual iter.Seq[V]) Checker {
	return asChecker(
		internal.Equal(slices.AppendSeq(make([]V, 0), expect)),
		slices.AppendSeq(make([]V, 0), actual),
	)
}

func NotEqual[V any](expect V, actual V) Checker {
	matcher := internal.NotEqual(expect)
	return asChecker(matcher, actual)
}

func MatchSnapshot(build func(s *snapshot.Snapshot), snapshotName string) Checker {
	return asChecker(snapshot.Match(snapshotName), Build(build))
}

func asChecker[T any](matcher internal.Matcher[T], actual T) Checker {
	return &checker[T]{
		Matcher: matcher,
		actual:  actual,
	}
}

type checker[T any] struct {
	internal.Matcher[T]

	actual T
}

func (c *checker[T]) Check(t TB) {
	switch x := t.(type) {
	case interface{ Unwrap() *testing.T }:
		tt := x.Unwrap()
		tt.Helper()
		internal.Expect(tt, c.actual, c.Matcher)
	}
}
