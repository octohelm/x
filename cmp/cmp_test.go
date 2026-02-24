package cmp_test

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/octohelm/x/cmp"
	. "github.com/octohelm/x/testing/v2"
)

type User struct {
	ID   int
	Name string
	Tags []string
}

func TestCmp(t *testing.T) {
	t.Run("原子比较 (True, False, Eq, Neq)", func(t *testing.T) {
		Then(t, "布尔校验",
			Expect(true, Be(cmp.True())),
			Expect(false, Be(cmp.False())),
		)

		Then(t, "等值校验",
			Expect(100, Be(cmp.Eq(100))),
			Expect("go", Be(cmp.Neq("java"))),
		)
	})

	t.Run("数值区间 (Gt, Gte, Lt, Lte)", func(t *testing.T) {
		val := 10
		Then(t, "区间断言",
			Expect(val, Be(cmp.Gt(5))),
			Expect(val, Be(cmp.Gte(10))),
			Expect(val, Be(cmp.Lt(20))),
			Expect(val, Be(cmp.Lte(10))),
		)
	})

	t.Run("状态校验 (Nil, NotNil, Zero, NotZero)", func(t *testing.T) {
		var ptr *int
		var s []int

		Then(t, "空指针与空容器",
			Expect(ptr, Be(cmp.Nil[*int]())),
			Expect(s, Be(cmp.Nil[[]int]())),
			Expect(1, Be(cmp.NotNil[int]())),
		)

		Then(t, "零值状态",
			Expect(0, Be(cmp.Zero[int]())),
			Expect(User{}, Be(cmp.Zero[User]())),
			Expect("hello", Be(cmp.NotZero[string]())),
		)
	})

	t.Run("容器长度 (Len)", func(t *testing.T) {
		list := []int{1, 2, 3}
		dict := map[string]int{"a": 1}

		Then(t, "支持 int 或谓词函数",
			Expect(list, Be(cmp.Len[[]int](3))),
			Expect(dict, Be(cmp.Len[map[string]int](cmp.Gt(0)))),
			Expect("golang", Be(cmp.Len[string](cmp.Lte(10)))),
		)
	})

	t.Run("迭代器校验 (Every)", func(t *testing.T) {
		t.Run("校验 Slice 元素", func(t *testing.T) {
			nums := []int{2, 4, 6, 8}
			// 配合 slices.Values 转换为 iter.Seq[int]
			Then(t, "所有元素必须为偶数",
				Expect(slices.Values(nums), Be(cmp.Every(func(v int) error {
					if v%2 != 0 {
						return fmt.Errorf("must be even")
					}
					return nil
				}))),
			)
		})

		t.Run("校验 Map 键值", func(t *testing.T) {
			m := map[string]int{"a": 10, "b": 20}

			Then(t, "校验所有 Value",
				Expect(maps.Values(m), Be(cmp.Every(cmp.Gte(10)))),
			)

			Then(t, "校验所有 Key",
				Expect(maps.Keys(m), Be(cmp.Every(cmp.NotZero[string]()))),
			)
		})
	})

	t.Run("迭代器校验 (Some)", func(t *testing.T) {
		nums := []int{1, 3, 5, 8, 9}

		Then(t, "存在满足条件的元素",
			// 集合里至少有一个偶数 (8)
			Expect(slices.Values(nums), Be(cmp.Some(func(v int) error {
				if v%2 == 0 {
					return nil
				}
				return fmt.Errorf("not even")
			}))),
		)

		Then(t, "配合内建谓词",
			Expect(slices.Values(nums), Be(cmp.Some(cmp.Eq(5)))),
		)
	})

	t.Run("综合复杂场景", func(t *testing.T) {
		users := []User{
			{ID: 1, Name: "Alice", Tags: []string{"admin", "staff"}},
			{ID: 2, Name: "Bob", Tags: []string{"staff"}},
		}

		Then(t, "深度嵌套校验",
			Expect(
				slices.Values(users),
				Be(cmp.Every(func(u User) error {
					if err := cmp.Gt(0)(u.ID); err != nil {
						return err
					}
					return cmp.Some(cmp.Eq("staff"))(
						slices.Values(u.Tags),
					)
				}))),
		)
	})

	t.Run("错误结构验证", func(t *testing.T) {
		nums := []int{1, 2, 3}
		err := cmp.Every(cmp.Eq(1))(slices.Values(nums))
		if err != nil {
			e, ok := errors.AsType[*cmp.ErrCheck](err)

			Then(t, "应返回 ErrCheck 类型且 Topic 为 elem",
				Expect(ok, Be(cmp.True())),
				Expect(e.Topic, Be(cmp.Eq("elem"))),
				Expect(e.Err, Be(cmp.NotNil[error]())),
			)
		}
	})
}
