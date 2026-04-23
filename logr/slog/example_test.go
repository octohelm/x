package slog_test

import (
	"context"
	"fmt"
	stdslog "log/slog"

	logrslog "github.com/octohelm/x/logr/slog"
)

func ExampleLogger() {
	logger := logrslog.Logger(logrslog.Default(logrslog.EnableLevel(stdslog.LevelInfo)))
	_, child := logger.Start(context.Background(), "sync", "key", "value")

	fmt.Printf("%T\n", child)
	// Output:
	// *slog.logger
}

func ExampleDefault() {
	logger := logrslog.Default(logrslog.EnableLevel(stdslog.LevelWarn))

	fmt.Println(logger.Enabled(context.Background(), stdslog.LevelInfo))
	fmt.Println(logger.Enabled(context.Background(), stdslog.LevelWarn))
	// Output:
	// false
	// true
}
