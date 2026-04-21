package logr_test

import (
	"context"
	"errors"
	"testing"

	"github.com/octohelm/x/cmp"
	"github.com/octohelm/x/logr"
	. "github.com/octohelm/x/testing/v2"
)

func TestLevel(t *testing.T) {
	t.Run("ParseLevel", func(t *testing.T) {
		cases := []struct {
			name  string
			input string
			level logr.Level
		}{
			{name: "error", input: "error", level: logr.ErrorLevel},
			{name: "warn", input: "warn", level: logr.WarnLevel},
			{name: "warning", input: "warning", level: logr.WarnLevel},
			{name: "info", input: "info", level: logr.InfoLevel},
			{name: "debug", input: "debug", level: logr.DebugLevel},
			{name: "mixed case", input: "DeBuG", level: logr.DebugLevel},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				level, err := logr.ParseLevel(c.input)

				Then(t, "should parse known level text",
					Expect(err, Equal[error](nil)),
					Expect(level, Equal(c.level)),
				)
			})
		}

		t.Run("invalid", func(t *testing.T) {
			_, err := logr.ParseLevel("verbose")

			Then(t, "unknown text should return error",
				Expect(err == nil, Be(cmp.False())),
			)
		})
	})

	t.Run("MarshalText and String", func(t *testing.T) {
		expected := map[logr.Level]string{
			logr.ErrorLevel: "error",
			logr.WarnLevel:  "warning",
			logr.InfoLevel:  "info",
			logr.DebugLevel: "debug",
		}

		for level, text := range expected {
			t.Run(text, func(t *testing.T) {
				got, err := level.MarshalText()

				Then(t, "marshal and string should agree",
					Expect(err, Equal[error](nil)),
					Expect(string(got), Equal(text)),
					Expect(level.String(), Equal(text)),
				)
			})
		}

		t.Run("unknown", func(t *testing.T) {
			level := logr.Level(99)
			_, err := level.MarshalText()

			Then(t, "unknown level should report fallback string",
				Expect(err == nil, Be(cmp.False())),
				Expect(level.String(), Equal("unknown")),
			)
		})
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		var level logr.Level
		err := level.UnmarshalText([]byte("warning"))

		Then(t, "should update receiver on success",
			Expect(err, Equal[error](nil)),
			Expect(level, Equal(logr.WarnLevel)),
		)

		err = level.UnmarshalText([]byte("verbose"))

		Then(t, "should keep parse error for unknown text",
			Expect(err == nil, Be(cmp.False())),
		)
	})
}

func TestDiscardLogger(t *testing.T) {
	t.Run("all methods are safe no-ops", func(t *testing.T) {
		logger := logr.Discard()
		child := logger.WithValues("k", "v")
		startedCtx, startedLogger := logger.Start(context.Background(), "span", "k", "v")

		logger.Debug("debug %d", 1)
		logger.Info("info %d", 1)
		logger.Warn(errors.New("warn"))
		logger.Error(errors.New("error"))
		logger.End()
		startedLogger.End()

		Then(t, "discard logger should stay usable across methods",
			Expect(logger == nil, Be(cmp.False())),
			Expect(child == nil, Be(cmp.False())),
			Expect(startedLogger == nil, Be(cmp.False())),
			Expect(startedCtx, Equal(context.Background())),
		)
	})
}
