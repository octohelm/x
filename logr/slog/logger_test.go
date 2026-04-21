package slog

import (
	"bytes"
	"context"
	stdslog "log/slog"
	"strings"
	"testing"
	"time"

	"github.com/octohelm/x/logr"
	. "github.com/octohelm/x/testing/v2"
)

func newJSONLogger(buf *bytes.Buffer, level stdslog.Level) *stdslog.Logger {
	return stdslog.New(stdslog.NewJSONHandler(buf, &stdslog.HandlerOptions{
		Level: level,
	}))
}

func TestHandlerAndLogger(t *testing.T) {
	t.Run("EnableLevel 控制 handler 输出级别", func(t *testing.T) {
		h := &handler{}
		EnableLevel(stdslog.LevelWarn)(h)

		Then(t, "低于阈值时禁用，高于阈值时启用",
			Expect(h.Enabled(context.Background(), stdslog.LevelInfo), Equal(false)),
			Expect(h.Enabled(context.Background(), stdslog.LevelWarn), Equal(true)),
			Expect(h.Enabled(context.Background(), stdslog.LevelError), Equal(true)),
		)
	})

	t.Run("Default 返回可用 logger", func(t *testing.T) {
		l := Default(EnableLevel(stdslog.LevelInfo))

		Then(t, "应成功构造 logger",
			Expect(l == nil, Equal(false)),
		)
	})

	t.Run("handler 透传 Handle、WithAttrs 和 WithGroup", func(t *testing.T) {
		buf := &bytes.Buffer{}
		base := stdslog.NewJSONHandler(buf, nil)
		h := &handler{h: base, lvl: stdslog.LevelDebug}

		record := stdslog.NewRecord(time.Now(), stdslog.LevelInfo, "hello", 0)
		record.AddAttrs(stdslog.String("k", "v"))

		withAttrs := h.WithAttrs([]stdslog.Attr{stdslog.String("scope", "test")})
		withGroup := withAttrs.WithGroup("worker")

		Then(t, "Handle 应透传到底层 handler",
			ExpectDo(func() error {
				return withGroup.Handle(context.Background(), record)
			}),
		)

		output := buf.String()
		Then(t, "输出应包含透传的 attrs 和 group",
			Expect(strings.Contains(output, "hello"), Equal(true)),
			Expect(strings.Contains(output, "scope"), Equal(true)),
			Expect(strings.Contains(output, "worker"), Equal(true)),
		)
	})

	t.Run("Logger 记录 info、warn 与 error", func(t *testing.T) {
		buf := &bytes.Buffer{}
		l := Logger(newJSONLogger(buf, stdslog.LevelInfo))

		l.Info("hello %s", "world")
		l.Warn(assertErr("warn"))
		l.Error(assertErr("error"))

		output := buf.String()

		Then(t, "应输出 info、warn 与 error 记录",
			Expect(strings.Contains(output, "hello world"), Equal(true)),
			Expect(strings.Contains(output, "\"level\":\"WARN\""), Equal(true)),
			Expect(strings.Contains(output, "\"level\":\"ERROR\""), Equal(true)),
		)
	})

	t.Run("Debug 会受级别控制", func(t *testing.T) {
		buf := &bytes.Buffer{}
		l := Logger(newJSONLogger(buf, stdslog.LevelInfo))

		l.Debug("debug %d", 1)

		Then(t, "info 级别下不应输出 debug",
			Expect(buf.String(), Equal("")),
		)
	})

	t.Run("WithValues 与 Start 生成结构化字段和 group", func(t *testing.T) {
		buf := &bytes.Buffer{}
		base := Logger(newJSONLogger(buf, stdslog.LevelDebug))

		ctx, child := base.Start(context.Background(), "sync-users", "worker", 1)
		withValues := child.WithValues("request_id", "r-1")
		withValues.Info("done")

		Then(t, "Start 返回原始上下文",
			Expect(ctx == context.Background(), Equal(true)),
		)

		output := buf.String()
		Then(t, "应输出 group 和附加字段",
			Expect(strings.Contains(output, "sync-users"), Equal(true)),
			Expect(strings.Contains(output, "request_id"), Equal(true)),
			Expect(strings.Contains(output, "worker"), Equal(true)),
		)
	})

	t.Run("End 会回退一个 span 层级", func(t *testing.T) {
		l := &logger{
			ctx:   context.Background(),
			slog:  newJSONLogger(&bytes.Buffer{}, stdslog.LevelDebug),
			spans: []string{"a", "b"},
		}

		l.End()

		Then(t, "End 后应移除最后一个 span",
			Expect(l.spans, Equal([]string{"a"})),
		)

		l.End()
		l.End()

		Then(t, "空 span 上再次 End 不应 panic 或继续下溢",
			Expect(l.spans, Equal([]string{})),
		)
	})

	t.Run("logr 上下文集成", func(t *testing.T) {
		buf := &bytes.Buffer{}
		ctx := logr.WithLogger(context.Background(), Logger(newJSONLogger(buf, stdslog.LevelInfo)))

		log := logr.FromContext(ctx)
		log.Info("context logger")

		Then(t, "从 context 读取的 logger 可正常输出",
			Expect(strings.Contains(buf.String(), "context logger"), Equal(true)),
		)
	})
}

type assertErr string

func (e assertErr) Error() string { return string(e) }
