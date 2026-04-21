package logr

import "context"

type contextKey struct{}

// WithLogger 将 Logger 写入上下文。
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

// FromContext 从上下文中读取 Logger；如果不存在则返回 Discard。
func FromContext(ctx context.Context) Logger {
	if v, ok := ctx.Value(contextKey{}).(Logger); ok {
		return v
	}
	return Discard()
}

// Start 使用上下文中的 Logger 开启一个新的日志作用域。
func Start(ctx context.Context, name string, keyAndValues ...any) (context.Context, Logger) {
	return FromContext(ctx).Start(ctx, name, keyAndValues...)
}

// LoggerFromContext 从上下文中读取 Logger，并报告是否存在。
func LoggerFromContext(ctx context.Context) (Logger, bool) {
	v, ok := ctx.Value(contextKey{}).(Logger)
	return v, ok
}

// LoggerInjectContext 将 Logger 注入上下文。
//
// 它与 WithLogger 等价，保留旧命名以兼容现有调用。
func LoggerInjectContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}
