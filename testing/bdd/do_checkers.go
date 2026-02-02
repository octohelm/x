package bdd

import (
	"testing"

	"github.com/octohelm/x/testing/internal"
)

func NotEqualDoValue[V any](expect V, do func() (V, error)) Checker {
	matcher := internal.NotEqual(expect)
	return asDoValueChecker(matcher, do)
}

func EqualDoValue[V any](expect V, do func() (V, error)) Checker {
	matcher := internal.Equal(expect)
	return asDoValueChecker(matcher, do)
}

func asDoValueChecker[T any](matcher internal.Matcher[T], do func() (T, error)) Checker {
	return &doValueChecker[T]{
		Matcher: matcher,
		do:      do,
	}
}

type doValueChecker[T any] struct {
	internal.Matcher[T]

	do func() (T, error)
}

func (c *doValueChecker[T]) Check(t TB) {
	switch x := t.(type) {
	case interface{ Unwrap() *testing.T }:
		tt := x.Unwrap()
		tt.Helper()

		actual, err := c.do()
		if err != nil {
			t.Fatal(err)
		}

		internal.Expect(tt, actual, c.Matcher)
	}
}
