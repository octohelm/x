package v2

import (
	"github.com/octohelm/x/testing/internal"
)

func Helper[X any](x X) X {
	return internal.Helper(1+1, x)
}

func Expect[V any](actual V, checkers ...ValueChecker[V]) Checker {
	return internal.Helper(1, &actionChecker[V]{
		do: func() (V, error) {
			return actual, nil
		},
		checkers: checkers,
	})
}

func ExpectMustValue[V any](do func() (V, error), checkers ...ValueChecker[V]) Checker {
	return internal.Helper(1, &actionChecker[V]{
		do:       do,
		checkers: checkers,
	})
}

type actionChecker[V any] struct {
	Reporter

	do       func() (V, error)
	checkers []ValueChecker[V]
}

func (r *actionChecker[V]) Check(t TB) {
	if len(r.checkers) == 0 {
		return
	}

	t.Helper()

	v, err := r.do()
	if err != nil {
		r.Fatal(t, err)
		return
	}

	for _, c := range r.checkers {
		c.Check(t, v)
	}
}

func ExpectMust(do func() error) Checker {
	return internal.Helper(1, &failureActionChecker{
		do: do,
	})
}

func ExpectDo(do func() error, errorCheckers ...ValueChecker[error]) Checker {
	return internal.Helper(1, &failureActionChecker{
		do:            do,
		errorCheckers: errorCheckers,
	})
}

type failureActionChecker struct {
	Reporter

	do            func() error
	errorCheckers []ValueChecker[error]
}

func (f *failureActionChecker) Check(t TB) {
	t.Helper()

	err := f.do()

	if len(f.errorCheckers) == 0 {
		// should be no error
		if err != nil {
			f.Fatal(t, err)
		}
		return
	}

	for _, errChecker := range f.errorCheckers {
		errChecker.Check(t, err)
	}
}
