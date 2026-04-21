package v2

import (
	"github.com/octohelm/x/testing/internal"
)

// Helper 标记调用栈层级，便于失败信息落在调用方位置。
func Helper[X any](x X) X {
	return internal.Helper(1+1, x)
}

// Expect 将实际值与一组 ValueChecker 组合为 Checker。
func Expect[V any](actual V, checkers ...ValueChecker[V]) Checker {
	return internal.Helper(1, &actionChecker[V]{
		do: func() (V, error) {
			return actual, nil
		},
		checkers: checkers,
	})
}

// ExpectMustValue 延迟执行取值动作，并在成功后应用一组 ValueChecker。
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

// ExpectMust 将返回 error 的动作包装为必须成功的 Checker。
func ExpectMust(do func() error) Checker {
	return internal.Helper(1, &failureActionChecker{
		do: do,
	})
}

// ExpectDo 将返回 error 的动作包装为 Checker，并允许继续断言该 error。
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
