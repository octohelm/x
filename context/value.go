package context

import (
	"context"
	"fmt"
	"reflect"
)

// WithValue 类似 context.WithValue，但跳过 key 的 comparable 检查。
//
// 它适合性能敏感且已自行保证 key 可安全比较的场景。
func WithValue(parent context.Context, key any, val any) context.Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if key == nil {
		panic("nil key")
	}
	return &valueCtx{parent, key, val}
}

type valueCtx struct {
	context.Context
	key, val any
}

func (c *valueCtx) String() string {
	return contextName(c.Context) + ".WithValue(type " + reflect.TypeOf(c.key).String() + ", val " + stringify(c.val) + ")"
}

func (c *valueCtx) Value(key any) any {
	if c.key == key {
		return c.val
	}
	return c.Context.Value(key)
}

func stringify(v any) string {
	switch s := v.(type) {
	case fmt.Stringer:
		return s.String()
	case string:
		return s
	}
	return "<not Stringer>"
}

func contextName(c context.Context) string {
	if s, ok := c.(fmt.Stringer); ok {
		return s.String()
	}
	return reflect.TypeOf(c).String()
}
