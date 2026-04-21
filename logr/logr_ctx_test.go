package logr_test

import (
	"context"
	"errors"
	"testing"

	"github.com/octohelm/x/logr"
	. "github.com/octohelm/x/testing/v2"
)

type stubLogger struct {
	lastStartName   string
	lastStartArgs   []any
	lastValues      []any
	debugMessages   []string
	infoMessages    []string
	warnErrors      []error
	errorErrors     []error
	endCalls        int
	returnChildSame bool
}

func (s *stubLogger) Start(ctx context.Context, name string, keyAndValues ...any) (context.Context, logr.Logger) {
	s.lastStartName = name
	s.lastStartArgs = append([]any(nil), keyAndValues...)
	if s.returnChildSame {
		return context.WithValue(ctx, "started", name), s
	}
	return context.WithValue(ctx, "started", name), &stubLogger{returnChildSame: true}
}

func (s *stubLogger) End() {
	s.endCalls++
}

func (s *stubLogger) WithValues(keyAndValues ...any) logr.Logger {
	next := *s
	next.lastValues = append([]any(nil), keyAndValues...)
	return &next
}

func (s *stubLogger) Debug(msg string, args ...any) {
	s.debugMessages = append(s.debugMessages, msg)
}

func (s *stubLogger) Info(msg string, args ...any) {
	s.infoMessages = append(s.infoMessages, msg)
}

func (s *stubLogger) Warn(err error) {
	s.warnErrors = append(s.warnErrors, err)
}

func (s *stubLogger) Error(err error) {
	s.errorErrors = append(s.errorErrors, err)
}

func TestContextHelpers(t *testing.T) {
	t.Run("缺省上下文返回 Discard logger", func(t *testing.T) {
		logger := logr.FromContext(context.Background())
		logger.Debug("debug")
		logger.Info("info")
		logger.Warn(errors.New("warn"))
		logger.Error(errors.New("error"))
		_, child := logger.Start(context.Background(), "span")
		child.End()

		Then(t, "缺省 logger 可安全调用所有方法",
			Expect(logger, Be(func(actual logr.Logger) error {
				if actual == nil {
					return errors.New("logger should not be nil")
				}
				return nil
			})),
		)
	})

	t.Run("WithLogger 与 LoggerFromContext", func(t *testing.T) {
		base := &stubLogger{}
		ctx := logr.WithLogger(context.Background(), base)

		logger, ok := logr.LoggerFromContext(ctx)

		Then(t, "应读回同一个 logger",
			Expect(ok, Equal(true)),
			Expect(logger, Equal[logr.Logger](base)),
			Expect(logr.FromContext(ctx), Equal[logr.Logger](base)),
		)
	})

	t.Run("LoggerInjectContext 与 WithLogger 等价", func(t *testing.T) {
		base := &stubLogger{}
		ctx := logr.LoggerInjectContext(context.Background(), base)

		Then(t, "注入后可通过 FromContext 读取",
			Expect(logr.FromContext(ctx), Equal[logr.Logger](base)),
		)
	})
}

func TestStartAndWithValues(t *testing.T) {
	t.Run("Start 透传到上下文中的 logger", func(t *testing.T) {
		base := &stubLogger{returnChildSame: true}
		ctx := logr.WithLogger(context.Background(), base)

		startedCtx, child := logr.Start(ctx, "sync-users", "worker", 1)

		Then(t, "Start 应记录名称与参数，并返回子 logger",
			Expect(base.lastStartName, Equal("sync-users")),
			Expect(base.lastStartArgs, Equal([]any{"worker", 1})),
			Expect(startedCtx.Value("started"), Equal(any("sync-users"))),
			Expect(child, Equal[logr.Logger](base)),
		)
	})

	t.Run("WithValues 返回绑定字段的新 logger", func(t *testing.T) {
		base := &stubLogger{}
		child := base.WithValues("k", "v").(*stubLogger)

		Then(t, "应保留新字段且不污染原 logger",
			Expect(child.lastValues, Equal([]any{"k", "v"})),
			Expect(base.lastValues, Equal([]any(nil))),
		)
	})
}
