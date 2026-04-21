package logr

import (
	"context"
)

// Logger 定义仓库内使用的统一日志接口。
type Logger interface {
	// Start 开启一个带名称的日志作用域，并返回新的上下文和子 Logger。
	Start(ctx context.Context, name string, keyAndValues ...any) (context.Context, Logger)
	// End 结束当前日志作用域。
	End()

	// WithValues 绑定结构化字段并返回新的 Logger。
	WithValues(keyAndValues ...any) Logger

	// Debug 记录调试级别日志。
	Debug(msg string, args ...any)
	// Info 记录信息级别日志。
	Info(msg string, args ...any)

	// Warn 记录警告级别错误。
	Warn(err error)

	// Error 记录错误级别错误。
	Error(err error)
}
