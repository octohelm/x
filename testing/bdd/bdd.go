package bdd

import (
	"context"
	"testing"
)

type TB interface {
	Chdir(dir string)
	Setenv(key, value string)

	Skip(args ...any)

	Context() context.Context
}

type T interface {
	TB

	Given(summary string, do func(b T))
	When(summary string, do func(b T))
	Then(summary string, checkers ...Checker)
}

func FromT(t *testing.T) T {
	return &bddT{T: t}
}

type bddT struct {
	*testing.T
}

func ScenarioT(setup func(b T)) func(*testing.T) {
	return func(t *testing.T) {
		setup(FromT(t))
	}
}

func GivenT(setup func(b T)) func(*testing.T) {
	return func(t *testing.T) {
		setup(FromT(t))
	}
}

func (t *bddT) Unwrap() *testing.T {
	return t.T
}

func (t *bddT) Given(summary string, setup func(b T)) {
	t.T.Run("GIVEN  "+summary, func(t *testing.T) {
		setup(FromT(t))
	})
}

func (t *bddT) When(summary string, setup func(b T)) {
	t.T.Run("WHEN  "+summary, func(t *testing.T) {
		setup(FromT(t))
	})
}

func (t *bddT) Then(summary string, checkers ...Checker) {
	t.T.Helper()

	t.T.Run("THEN  "+summary, func(t *testing.T) {
		t.Helper()

		tt := FromT(t)

		for _, c := range checkers {
			if c != nil {
				c.Check(tt)
			}
		}
	})
}
