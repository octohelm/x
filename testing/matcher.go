package testing

import (
	"github.com/octohelm/x/testing/internal"
)

type Matcher[A any] = internal.Matcher[A]

type ExpectedFormatter = internal.ExpectedFormatter

func NewMatcher[A any](name string, match func(a A) bool) Matcher[A] {
	return internal.NewMatcher(name, match)
}

func NewCompareMatcher[A any, E any](name string, match func(a A, e E) bool) func(e E) Matcher[A] {
	return func(expected E) Matcher[A] {
		return internal.NewCompareMatcher(name, match)(expected)
	}
}

func Not[A any](matcher Matcher[A]) Matcher[A] {
	return internal.Not(matcher)
}
