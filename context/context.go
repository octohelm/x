package context

import (
	"context"

	"github.com/pkg/errors"
)

func WithDefaultsFunc[T any](defaultsFunc func() T) OptionFunc[T] {
	return func(c *ctx[T]) {
		c.defaultsFunc = defaultsFunc
	}
}

func WithDefaults[T any](v T) OptionFunc[T] {
	return func(c *ctx[T]) {
		c.defaultsFunc = func() T {
			return v
		}
	}
}

type OptionFunc[T any] func(c *ctx[T])

func New[T any](optFns ...OptionFunc[T]) Context[T] {
	c := &ctx[T]{}
	for _, fn := range optFns {
		fn(c)
	}
	return c
}

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
	panic(errors.Errorf("%T not found in context", c))
}

func (c *ctx[T]) MayFrom(ctx context.Context) (T, bool) {
	if v, ok := ctx.Value(c).(T); ok {
		return v, true
	}
	return *new(T), false
}
