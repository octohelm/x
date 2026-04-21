package context

import (
	stdcontext "context"
	"fmt"
	"testing"

	"github.com/octohelm/x/cmp"
	. "github.com/octohelm/x/testing/v2"
)

type stringKey struct{}

type stringValue struct{ value string }

func (v stringValue) String() string {
	return v.value
}

func TestTypedContext(t *testing.T) {
	t.Run("可注入并取回值", func(t *testing.T) {
		slot := New[string]()
		ctx := slot.Inject(stdcontext.Background(), "v1")

		v, ok := slot.MayFrom(ctx)

		Then(t, "注入后的值可被读取",
			Expect(slot.From(ctx), Equal("v1")),
			Expect(v, Equal("v1")),
			Expect(ok, Be(cmp.True())),
		)
	})

	t.Run("缺失值时可返回默认值", func(t *testing.T) {
		calls := 0
		slot := New[int](WithDefaultsFunc(func() int {
			calls++
			return 42
		}))

		Then(t, "From 使用延迟默认值而 MayFrom 仍报告缺失",
			Expect(slot.From(stdcontext.Background()), Equal(42)),
			Expect(slot.From(stdcontext.Background()), Equal(42)),
		)

		v, ok := slot.MayFrom(stdcontext.Background())

		Then(t, "默认值函数按需执行且不影响 MayFrom",
			Expect(calls, Equal(2)),
			Expect(v, Equal(0)),
			Expect(ok, Be(cmp.False())),
		)
	})

	t.Run("固定默认值", func(t *testing.T) {
		slot := New[string](WithDefaults("fallback"))

		Then(t, "未注入时返回固定默认值",
			Expect(slot.From(stdcontext.Background()), Equal("fallback")),
		)
	})

	t.Run("缺失值且无默认值时会 panic", func(t *testing.T) {
		slot := New[string]()
		msg := panicMessage(func() {
			slot.From(stdcontext.Background())
		})

		Then(t, "panic 信息包含类型化上下文",
			Expect(msg, Equal("*context.ctx[string] not found in context")),
		)
	})
}

func TestWithValue(t *testing.T) {
	t.Run("读取链路保持兼容", func(t *testing.T) {
		type key struct{ name string }

		parent := stdcontext.WithValue(stdcontext.Background(), key{name: "parent"}, "p")
		ctx := WithValue(parent, key{name: "child"}, "c")

		Then(t, "可读取当前值与父级值",
			Expect(ctx.Value(key{name: "child"}), Equal(any("c"))),
			Expect(ctx.Value(key{name: "parent"}), Equal(any("p"))),
		)
	})

	t.Run("非法输入会 panic", func(t *testing.T) {
		nilParentMsg := panicMessage(func() {
			WithValue(nil, "k", "v")
		})
		nilKeyMsg := panicMessage(func() {
			WithValue(stdcontext.Background(), nil, "v")
		})

		Then(t, "nil parent 与 nil key 会立即失败",
			Expect(nilParentMsg, Equal("cannot create context from nil parent")),
			Expect(nilKeyMsg, Equal("nil key")),
		)
	})

	t.Run("字符串表示优先使用 Stringer", func(t *testing.T) {
		ctx := WithValue(stdcontext.Background(), stringKey{}, stringValue{value: "demo"})

		Then(t, "String 输出包含 key 类型与格式化后的值",
			Expect(ctx.(*valueCtx).String(), Equal("context.Background.WithValue(type context.stringKey, val demo)")),
		)
	})
}

func panicMessage(do func()) (msg string) {
	defer func() {
		if r := recover(); r != nil {
			msg = fmt.Sprint(r)
		}
	}()
	do()
	return ""
}
