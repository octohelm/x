package iter_test

import (
	"errors"
	"testing"

	xiter "github.com/octohelm/x/iter"
	. "github.com/octohelm/x/testing/v2"
)

func TestAction(t *testing.T) {
	t.Run("正常产出多个值", func(t *testing.T) {
		seq := xiter.Action(func(yield func(*int) bool) error {
			for _, v := range []int{1, 2, 3} {
				v := v
				if !yield(&v) {
					return nil
				}
			}
			return nil
		})

		values := make([]int, 0)
		errs := make([]error, 0)

		for v, err := range seq {
			if v != nil {
				values = append(values, *v)
			}
			if err != nil {
				errs = append(errs, err)
			}
		}

		Then(t, "应该按顺序返回所有值且没有错误",
			Expect(values, Equal([]int{1, 2, 3})),
			Expect(errs, Equal([]error{})),
		)
	})

	t.Run("消费者提前终止", func(t *testing.T) {
		calls := 0

		seq := xiter.Action(func(yield func(*int) bool) error {
			for _, v := range []int{1, 2, 3} {
				v := v
				calls++
				if !yield(&v) {
					return nil
				}
			}
			return nil
		})

		values := make([]int, 0, 1)

		for v, err := range seq {
			_ = err
			values = append(values, *v)
			break
		}

		Then(t, "应该只消费首个值并及时停止上游 yield",
			Expect(values, Equal([]int{1})),
			Expect(calls, Equal(1)),
		)
	})

	t.Run("上游返回错误", func(t *testing.T) {
		expectErr := errors.New("boom")

		seq := xiter.Action(func(yield func(*int) bool) error {
			v := 1
			if !yield(&v) {
				return nil
			}
			return expectErr
		})

		values := make([]int, 0, 1)
		var gotErr error

		for v, err := range seq {
			if v != nil {
				values = append(values, *v)
			}
			if err != nil {
				gotErr = err
			}
		}

		Then(t, "应该先返回值，再在末尾返回错误",
			Expect(values, Equal([]int{1})),
			Expect(gotErr, Equal(expectErr)),
		)
	})
}
