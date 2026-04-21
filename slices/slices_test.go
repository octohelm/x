package slices_test

import (
	"strconv"
	"testing"

	xslices "github.com/octohelm/x/slices"
	. "github.com/octohelm/x/testing/v2"
)

func TestMap(t *testing.T) {
	t.Run("保持输入顺序映射", func(t *testing.T) {
		got := xslices.Map([]int{1, 2, 3}, func(v int) string {
			return strconv.Itoa(v * 10)
		})

		Then(t, "应按原顺序返回映射结果",
			Expect(got, Equal([]string{"10", "20", "30"})),
		)
	})

	t.Run("nil 切片输入", func(t *testing.T) {
		var input []int

		got := xslices.Map(input, func(v int) int {
			return v + 1
		})

		Then(t, "应返回空结果切片",
			Expect(got, Equal([]int{})),
		)
	})
}

func TestFilter(t *testing.T) {
	t.Run("按条件筛选且保持顺序", func(t *testing.T) {
		got := xslices.Filter([]int{1, 2, 3, 4}, func(v int) bool {
			return v%2 == 0
		})

		Then(t, "应仅保留满足条件的元素",
			Expect(got, Equal([]int{2, 4})),
		)
	})

	t.Run("无元素满足条件", func(t *testing.T) {
		got := xslices.Filter([]int{1, 3, 5}, func(v int) bool {
			return v%2 == 0
		})

		Then(t, "应返回空切片而不是 nil",
			Expect(got, Equal([]int{})),
		)
	})

	t.Run("nil 切片输入", func(t *testing.T) {
		var input []int

		got := xslices.Filter(input, func(v int) bool {
			return v > 0
		})

		Then(t, "应返回空结果切片",
			Expect(got, Equal([]int{})),
		)
	})
}
