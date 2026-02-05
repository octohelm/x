package testing

import (
	"github.com/octohelm/x/testing/internal/deprecated"
)

type Matcher[A any] = deprecated.Matcher[A]

type MatcherWithNormalizedExpected = deprecated.MatcherWithNormalizedExpected

func NewMatcher[A any](name string, match func(a A) bool) Matcher[A] {
	return deprecated.NewMatcher(name, match)
}

func NewCompareMatcher[A any, E any](name string, match func(a A, e E) bool) func(e E) Matcher[A] {
	return func(expected E) Matcher[A] {
		return deprecated.NewCompareMatcher(name, match)(expected)
	}
}

func Not[A any](matcher Matcher[A]) Matcher[A] {
	return deprecated.Not(matcher)
}
