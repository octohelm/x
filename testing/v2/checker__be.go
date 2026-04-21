package v2

import "github.com/octohelm/x/testing/internal"

// Be 将返回 error 的谓词包装为 ValueChecker。
func Be[V any](v func(v V) error) ValueChecker[V] {
	return internal.Helper(1, &beChecker[V]{
		be: v,
	})
}

type beChecker[V any] struct {
	Reporter

	be func(v V) error
}

func (r *beChecker[V]) Check(t TB, actual V) {
	t.Helper()

	if err := r.be(actual); err != nil {
		r.Fatal(t, err)
	}
}
