package cmp_test

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"testing"

	"github.com/go-json-experiment/json/jsontext"

	"github.com/octohelm/x/cmp"
	. "github.com/octohelm/x/testing/v2"
)

type user struct {
	ID   int
	Name string
	Tags []string
}

func TestAtomicPredicates(t *testing.T) {
	t.Run("成功路径", func(t *testing.T) {
		Then(t, "原子谓词可通过基本断言",
			Expect(cmp.True()(true), Equal[error](nil)),
			Expect(cmp.False()(false), Equal[error](nil)),
			Expect(cmp.Eq(100)(100), Equal[error](nil)),
			Expect(cmp.Neq("go")("java"), Equal[error](nil)),
			Expect(cmp.Gt(5)(10), Equal[error](nil)),
			Expect(cmp.Gte(10)(10), Equal[error](nil)),
			Expect(cmp.Lt(20)(10), Equal[error](nil)),
			Expect(cmp.Lte(10)(10), Equal[error](nil)),
		)
	})

	t.Run("失败路径", func(t *testing.T) {
		errTrue := cmp.True()(false)
		errNeq := cmp.Neq("go")("go")

		var condErr *cmp.ErrCondition

		Then(t, "失败时返回条件错误并保留操作符",
			Expect(errTrue, ErrorAsType[*cmp.ErrCondition]()),
			Expect(errors.As(errTrue, &condErr), Be(cmp.True())),
			Expect(condErr.Op, Equal("==")),
			Expect(errNeq, ErrorAsType[*cmp.ErrCondition]()),
		)
	})
}

func TestStatePredicates(t *testing.T) {
	var nilPtr *int
	var nilSlice []int
	var nilErr error

	t.Run("空值与零值成功路径", func(t *testing.T) {
		Then(t, "可识别 nil 与零值",
			Expect(cmp.Nil[*int]()(nilPtr), Equal[error](nil)),
			Expect(cmp.Nil[[]int]()(nilSlice), Equal[error](nil)),
			Expect(cmp.Nil[error]()(nilErr), Equal[error](nil)),
			Expect(cmp.NotNil[int]()(1), Equal[error](nil)),
			Expect(cmp.Zero[int]()(0), Equal[error](nil)),
			Expect(cmp.Zero[user]()(user{}), Equal[error](nil)),
			Expect(cmp.NotZero[string]()("hello"), Equal[error](nil)),
		)
	})

	t.Run("失败路径", func(t *testing.T) {
		errNotNil := cmp.NotNil[*int]()(nilPtr)
		errZero := cmp.Zero[string]()("hello")

		Then(t, "失败时返回状态错误",
			Expect(errNotNil, ErrorAsType[*cmp.ErrState]()),
			Expect(errZero, ErrorAsType[*cmp.ErrState]()),
		)
	})
}

func TestLen(t *testing.T) {
	t.Run("支持长度比较与谓词", func(t *testing.T) {
		Then(t, "容器长度可直接比较或继续套用谓词",
			Expect(cmp.Len[[]int](3)([]int{1, 2, 3}), Equal[error](nil)),
			Expect(cmp.Len[map[string]int](cmp.Gt(0))(map[string]int{"a": 1}), Equal[error](nil)),
			Expect(cmp.Len[string](cmp.Lte(10))("golang"), Equal[error](nil)),
		)
	})

	t.Run("非法类型", func(t *testing.T) {
		err := cmp.Len[int](1)(10)

		Then(t, "非长度类型返回状态错误",
			Expect(err, ErrorAsType[*cmp.ErrState]()),
		)
	})

	t.Run("嵌套谓词失败", func(t *testing.T) {
		err := cmp.Len[string](cmp.Gt(3))("go")

		var checkErr *cmp.ErrCheck

		Then(t, "长度失败会包装为 len 主题错误",
			Expect(err, ErrorAsType[*cmp.ErrCheck]()),
			Expect(errors.As(err, &checkErr), Be(cmp.True())),
			Expect(checkErr.Topic, Equal("len")),
			Expect(checkErr.Err, ErrorAsType[*cmp.ErrCondition]()),
		)
	})
}

func TestIterators(t *testing.T) {
	t.Run("Every 成功路径", func(t *testing.T) {
		nums := []int{2, 4, 6, 8}
		m := map[string]int{"a": 10, "b": 20}

		Then(t, "每个元素都满足谓词",
			Expect(cmp.Every(func(v int) error {
				if v%2 != 0 {
					return fmt.Errorf("must be even")
				}
				return nil
			})(slices.Values(nums)), Equal[error](nil)),
			Expect(cmp.Every(cmp.Gte(10))(maps.Values(m)), Equal[error](nil)),
			Expect(cmp.Every(cmp.NotZero[string]())(maps.Keys(m)), Equal[error](nil)),
		)
	})

	t.Run("Every 失败路径会保留索引指针", func(t *testing.T) {
		users := []user{
			{ID: 1, Name: "Alice", Tags: []string{"admin", "staff"}},
			{ID: 2, Name: "Bob", Tags: []string{}},
		}

		err := cmp.Every(func(u user) error {
			if err := cmp.Gt(0)(u.ID); err != nil {
				return err
			}
			return cmp.Len[[]string](cmp.Gt(0))(u.Tags)
		})(slices.Values(users))

		var checkErr *cmp.ErrCheck

		Then(t, "嵌套失败会保留元素索引和内部主题",
			Expect(err, ErrorAsType[*cmp.ErrCheck]()),
			Expect(errors.As(err, &checkErr), Be(cmp.True())),
			Expect(checkErr.Pointer, Equal(jsontext.Pointer("/1"))),
			Expect(checkErr.Topic, Equal("len")),
		)
	})

	t.Run("Some 的成功与失败路径", func(t *testing.T) {
		nums := []int{1, 3, 5, 8, 9}
		err := cmp.Some(cmp.Eq(2))(slices.Values([]int{1, 3, 5}))

		var checkErr *cmp.ErrCheck

		Then(t, "至少一个元素满足时返回成功",
			Expect(cmp.Some(cmp.Eq(5))(slices.Values(nums)), Equal[error](nil)),
		)

		Then(t, "无元素满足时返回 elem 主题错误",
			Expect(err, ErrorAsType[*cmp.ErrCheck]()),
			Expect(errors.As(err, &checkErr), Be(cmp.True())),
			Expect(checkErr.Topic, Equal("elem")),
			Expect(checkErr.Err, ErrorMatch(regexp.MustCompile(`none of the elements satisfy the predicate`))),
		)
	})
}
