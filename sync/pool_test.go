package sync_test

import (
	"testing"

	xsync "github.com/octohelm/x/sync"
	. "github.com/octohelm/x/testing/v2"
)

func TestPool(t *testing.T) {
	t.Run("零值池返回类型零值", func(t *testing.T) {
		var p xsync.Pool[int]

		Then(t, "空池且未设置 New 时返回零值",
			Expect(p.Get(), Equal(0)),
		)
	})

	t.Run("设置 New 后按需创建对象", func(t *testing.T) {
		calls := 0
		p := xsync.Pool[int]{
			New: func() int {
				calls++
				return 42
			},
		}

		first := p.Get()
		second := p.Get()

		Then(t, "Get 应通过 New 懒创建对象",
			Expect(first, Equal(42)),
			Expect(second, Equal(42)),
			Expect(calls, Equal(2)),
		)
	})

	t.Run("放回对象后优先复用池中值", func(t *testing.T) {
		p := xsync.Pool[string]{
			New: func() string {
				return "new"
			},
		}

		p.Put("cached")

		Then(t, "池中已有对象时应优先返回缓存值",
			Expect(p.Get(), Equal("cached")),
			Expect(p.Get(), Equal("new")),
		)
	})
}
