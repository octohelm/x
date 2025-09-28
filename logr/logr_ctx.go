package logr

import "context"

type contextKey struct{}

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

func FromContext(ctx context.Context) Logger {
	if v, ok := ctx.Value(contextKey{}).(Logger); ok {
		return v
	}
	return Discard()
}

func Start(ctx context.Context, name string, keyAndValues ...any) (context.Context, Logger) {
	return FromContext(ctx).Start(ctx, name, keyAndValues...)
}

func LoggerFromContext(ctx context.Context) (Logger, bool) {
	v, ok := ctx.Value(contextKey{}).(Logger)
	return v, ok
}

func LoggerInjectContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}
