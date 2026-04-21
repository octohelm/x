package context

import (
	"context"
	"fmt"
)

// WithDefaultsFunc 为 Context[T] 配置延迟求值的默认值函数。
func WithDefaultsFunc[T any](defaultsFunc func() T) OptionFunc[T] {
	return func(c *ctx[T]) {
		c.defaultsFunc = defaultsFunc
	}
}

// WithDefaults 为 Context[T] 配置固定默认值。
func WithDefaults[T any](v T) OptionFunc[T] {
	return func(c *ctx[T]) {
		c.defaultsFunc = func() T {
			return v
		}
	}
}

// OptionFunc 表示创建类型化 Context 时的配置项。
type OptionFunc[T any] func(c *ctx[T])

// New 创建一个类型化的上下文槽位。
func New[T any](optFns ...OptionFunc[T]) Context[T] {
	c := &ctx[T]{}
	for _, fn := range optFns {
		fn(c)
	}
	return c
}

// Context 表示一个可向 context.Context 注入和读取特定类型值的槽位。
type Context[T any] interface {
	Inject(ctx context.Context, value T) context.Context
	From(ctx context.Context) T
	MayFrom(ctx context.Context) (T, bool)
}

type ctx[T any] struct {
	defaultsFunc func() T
}

func (c *ctx[T]) Inject(ctx context.Context, value T) context.Context {
	return WithValue(ctx, c, value)
}

func (c *ctx[T]) From(ctx context.Context) T {
	if v, ok := ctx.Value(c).(T); ok {
		return v
	}
	if c.defaultsFunc != nil {
		return c.defaultsFunc()
	}
	panic(fmt.Errorf("%T not found in context", c))
}

func (c *ctx[T]) MayFrom(ctx context.Context) (T, bool) {
	if v, ok := ctx.Value(c).(T); ok {
		return v, true
	}
	return *new(T), false
}
